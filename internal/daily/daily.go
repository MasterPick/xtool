// Package daily 日常工具模块
// 提供计算器、单位换算、汇率换算、备忘录等日常实用工具
package daily

import (
	"crypto/rand"
	"fmt"
	"html"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"xtool/internal/db"
)

// DailyTools 日常工具结构体（Wails 绑定到前端）
type DailyTools struct {
	db *db.Database // 数据库连接
}

// Note 备忘录结构
type Note struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Color     string `json:"color"`
	Pinned    bool   `json:"pinned"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ConversionResult 单位转换结果
type ConversionResult struct {
	Value   float64 `json:"value"`   // 转换后的值
	Unit    string  `json:"unit"`    // 目标单位
	Formula string  `json:"formula"` // 转换公式说明
}

// NewDailyTools 创建日常工具模块实例
func NewDailyTools(database *db.Database) *DailyTools {
	return &DailyTools{db: database}
}

// ============================================================
// 计算器（标准模式 + 科学模式）
// ============================================================

// CalcBasic 基础四则运算
// op: +, -, *, /
func (d *DailyTools) CalcBasic(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("不支持的运算符: %s", op)
	}
}

// CalcScientific 科学计算函数
// func: sqrt/pow/log/ln/sin/cos/tan/abs/ceil/floor/round
func (d *DailyTools) CalcScientific(value, param float64, fn string) (float64, error) {
	switch fn {
	case "sqrt":
		if value < 0 {
			return 0, fmt.Errorf("负数不能开平方根")
		}
		return math.Sqrt(value), nil
	case "pow":
		return math.Pow(value, param), nil
	case "log":
		if value <= 0 {
			return 0, fmt.Errorf("对数的真数必须大于 0")
		}
		return math.Log10(value), nil
	case "ln":
		if value <= 0 {
			return 0, fmt.Errorf("自然对数的真数必须大于 0")
		}
		return math.Log(value), nil
	case "sin":
		return math.Sin(value * math.Pi / 180), nil // 角度转弧度
	case "cos":
		return math.Cos(value * math.Pi / 180), nil
	case "tan":
		// tan 函数在 90 度和 270 度附近趋向无穷大，添加边界处理
		// 将角度归一化到 0-360 范围
		normalized := math.Mod(value, 360)
		if normalized < 0 {
			normalized += 360
		}
		// 检查是否接近 90 度或 270 度（奇数倍的 90 度）
		if math.Mod(normalized, 180) == 90 {
			return 0, fmt.Errorf("tan(%.1f 度) 无定义（趋向无穷大）", value)
		}
		// 接近 90 度或 270 度时，检查是否会导致极大值
		radians := normalized * math.Pi / 180
		result := math.Tan(radians)
		if math.IsInf(result, 0) || math.IsNaN(result) {
			return 0, fmt.Errorf("tan(%.1f 度) 无定义（趋向无穷大）", value)
		}
		return result, nil
	case "abs":
		return math.Abs(value), nil
	case "ceil":
		return math.Ceil(value), nil
	case "floor":
		return math.Floor(value), nil
	case "round":
		return math.Round(value), nil
	case "pi":
		return math.Pi * value, nil
	default:
		return 0, fmt.Errorf("不支持的函数: %s", fn)
	}
}

// CalcExpression 安全表达式求值（替代前端 Function()）
// 支持基本四则运算、括号、数学常量（PI、E）
// 不使用 eval/exec，通过安全的词法分析实现
func (d *DailyTools) CalcExpression(expr string) (float64, error) {
	// 预处理：去除空格
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return 0, fmt.Errorf("表达式不能为空")
	}

	// 替换数学常量
	expr = strings.ReplaceAll(expr, "PI", fmt.Sprintf("%v", math.Pi))
	expr = strings.ReplaceAll(expr, "E", fmt.Sprintf("%v", math.E))

	// 安全检查：只允许数字、运算符、括号和小数点
	for _, ch := range expr {
		if !unicode.IsDigit(ch) && ch != '+' && ch != '-' && ch != '*' && ch != '/' &&
			ch != '(' && ch != ')' && ch != '.' && ch != ' ' {
			return 0, fmt.Errorf("表达式包含非法字符: %c（仅支持数字和 +-*/()）", ch)
		}
	}

	// 使用安全的方式解析表达式
	result, err := safeEval(expr)
	if err != nil {
		return 0, fmt.Errorf("表达式计算失败: %v", err)
	}
	return result, nil
}

// safeEval 安全的表达式求值器
// 支持加减乘除和括号，使用递归下降解析
func safeEval(expr string) (float64, error) {
	// 词法分析：将表达式拆分为 token
	tokens := tokenize(expr)
	pos := 0

	result, err := parseAddSub(tokens, &pos)
	if err != nil {
		return 0, err
	}
	if pos != len(tokens) {
		return 0, fmt.Errorf("表达式解析不完整，剩余部分: %s", strings.Join(tokens[pos:], " "))
	}
	return result, nil
}

// parseAddSub 解析加减法（最低优先级）
func parseAddSub(tokens []string, pos *int) (float64, error) {
	left, err := parseMulDiv(tokens, pos)
	if err != nil {
		return 0, err
	}
	for *pos < len(tokens) && (tokens[*pos] == "+" || tokens[*pos] == "-") {
		op := tokens[*pos]
		*pos++
		right, err := parseMulDiv(tokens, pos)
		if err != nil {
			return 0, err
		}
		if op == "+" {
			left += right
		} else {
			left -= right
		}
	}
	return left, nil
}

// parseMulDiv 解析乘除法
func parseMulDiv(tokens []string, pos *int) (float64, error) {
	left, err := parsePrimary(tokens, pos)
	if err != nil {
		return 0, err
	}
	for *pos < len(tokens) && (tokens[*pos] == "*" || tokens[*pos] == "/") {
		op := tokens[*pos]
		*pos++
		right, err := parsePrimary(tokens, pos)
		if err != nil {
			return 0, err
		}
		if op == "*" {
			left *= right
		} else {
			if right == 0 {
				return 0, fmt.Errorf("除数不能为零")
			}
			left /= right
		}
	}
	return left, nil
}

// parsePrimary 解析基本元素（数字或括号表达式）
func parsePrimary(tokens []string, pos *int) (float64, error) {
	if *pos >= len(tokens) {
		return 0, fmt.Errorf("表达式不完整")
	}

	token := tokens[*pos]

	// 处理负号（一元运算符）
	if token == "-" {
		*pos++
		val, err := parsePrimary(tokens, pos)
		if err != nil {
			return 0, err
		}
		return -val, nil
	}

	// 处理正号（一元运算符）
	if token == "+" {
		*pos++
		return parsePrimary(tokens, pos)
	}

	// 处理括号
	if token == "(" {
		*pos++ // 跳过 '('
		val, err := parseAddSub(tokens, pos)
		if err != nil {
			return 0, err
		}
		if *pos >= len(tokens) || tokens[*pos] != ")" {
			return 0, fmt.Errorf("缺少右括号")
		}
		*pos++ // 跳过 ')'
		return val, nil
	}

	// 处理数字
	val, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0, fmt.Errorf("无法解析数字: %s", token)
	}
	*pos++
	return val, nil
}

// tokenize 将表达式字符串拆分为 token 列表
func tokenize(expr string) []string {
	var tokens []string
	var current strings.Builder

	for _, ch := range expr {
		if ch == ' ' {
			continue
		}
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '(' || ch == ')' {
			// 如果之前有累积的数字，先添加
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(ch))
		} else {
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// ============================================================
// 单位换算工具
// ============================================================

// ConvertLength 长度单位换算
// fromUnit/toUnit: mm, cm, m, km, inch, foot, yard, mile
func (d *DailyTools) ConvertLength(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为米（基准单位）
	toMeter := map[string]float64{
		"mm":   0.001,
		"cm":   0.01,
		"m":    1.0,
		"km":   1000.0,
		"inch": 0.0254,
		"foot": 0.3048,
		"yard": 0.9144,
		"mile": 1609.344,
	}

	fromFactor, ok1 := toMeter[fromUnit]
	toFactor, ok2 := toMeter[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的长度单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertWeight 重量单位换算
// fromUnit/toUnit: mg, g, kg, t, oz, lb
func (d *DailyTools) ConvertWeight(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为克（基准单位）
	toGram := map[string]float64{
		"mg": 0.001,
		"g":  1.0,
		"kg": 1000.0,
		"t":  1000000.0,
		"oz": 28.3495,
		"lb": 453.592,
	}

	fromFactor, ok1 := toGram[fromUnit]
	toFactor, ok2 := toGram[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的重量单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertTemperature 温度换算
// fromUnit/toUnit: C（摄氏度）、F（华氏度）、K（开尔文）
func (d *DailyTools) ConvertTemperature(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 先转换为摄氏度
	var celsius float64
	switch strings.ToUpper(fromUnit) {
	case "C":
		celsius = value
	case "F":
		celsius = (value - 32) * 5 / 9
	case "K":
		celsius = value - 273.15
	default:
		return ConversionResult{}, fmt.Errorf("不支持的温度单位（支持: C/F/K）")
	}

	// 再从摄氏度转为目标单位
	var result float64
	switch strings.ToUpper(toUnit) {
	case "C":
		result = celsius
	case "F":
		result = celsius*9/5 + 32
	case "K":
		result = celsius + 273.15
	default:
		return ConversionResult{}, fmt.Errorf("不支持的温度单位（支持: C/F/K）")
	}

	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g%s = %g%s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertSpeed 速度单位换算
// fromUnit/toUnit: ms（米/秒）、kmh（千米/时）、mph（英里/时）、knot（节）
func (d *DailyTools) ConvertSpeed(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为米/秒
	toMPS := map[string]float64{
		"ms":   1.0,
		"kmh":  1.0 / 3.6,
		"mph":  0.44704,
		"knot": 0.514444,
	}

	fromFactor, ok1 := toMPS[fromUnit]
	toFactor, ok2 := toMPS[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的速度单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertArea 面积单位换算
// fromUnit/toUnit: mm2, cm2, m2, km2, ha（公顷）, acre（英亩）, mu（亩）
func (d *DailyTools) ConvertArea(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为平方米（基准单位）
	toSqm := map[string]float64{
		"mm2":  1e-6,
		"cm2":  1e-4,
		"m2":   1.0,
		"km2":  1e6,
		"ha":   1e4,      // 公顷
		"acre": 4046.8564, // 英亩
		"mu":   666.6667,  // 亩
	}

	fromFactor, ok1 := toSqm[fromUnit]
	toFactor, ok2 := toSqm[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的面 积单位（支持: mm2/cm2/m2/km2/ha/acre/mu）")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertVolume 体积单位换算
// fromUnit/toUnit: ml, l（升）, m3, cm3, gallon（加仑）, oz_fl（液量盎司）
func (d *DailyTools) ConvertVolume(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为升（基准单位）
	toLiter := map[string]float64{
		"ml":    0.001,
		"l":     1.0,
		"m3":    1000.0,
		"cm3":   0.001,
		"gallon": 3.78541, // 美制加仑
		"oz_fl": 0.0295735, // 液量盎司
	}

	fromFactor, ok1 := toLiter[fromUnit]
	toFactor, ok2 := toLiter[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的体积单位（支持: ml/l/m3/cm3/gallon/oz_fl）")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertDataSize 数据存储单位换算
// fromUnit/toUnit: b（字节）, kb, mb, gb, tb, pb, bit（位）
func (d *DailyTools) ConvertDataSize(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为字节（基准单位）
	toByte := map[string]float64{
		"bit": 0.125,  // 1 字节 = 8 位
		"b":   1.0,
		"kb":  1024.0,
		"mb":  1024 * 1024,
		"gb":  1024 * 1024 * 1024,
		"tb":  1024.0 * 1024 * 1024 * 1024,
		"pb":  1024.0 * 1024 * 1024 * 1024 * 1024,
	}

	fromFactor, ok1 := toByte[fromUnit]
	toFactor, ok2 := toByte[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的数据存储单位（支持: bit/b/kb/mb/gb/tb/pb）")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    strings.ToUpper(toUnit),
		Formula: fmt.Sprintf("%g %s = %g %s", value, strings.ToUpper(fromUnit), result, strings.ToUpper(toUnit)),
	}, nil
}

// ============================================================
// 时间日期工具
// ============================================================

// GetCurrentTime 获取当前时间的多种格式表示
func (d *DailyTools) GetCurrentTime() map[string]string {
	now := time.Now()
	return map[string]string{
		"datetime":  now.Format("2006-01-02 15:04:05"),
		"date":      now.Format("2006-01-02"),
		"time":      now.Format("15:04:05"),
		"weekday":   now.Weekday().String(),
		"timestamp": fmt.Sprintf("%d", now.Unix()),
		"utc":       now.UTC().Format("2006-01-02 15:04:05 UTC"),
		"timezone":  now.Format("MST"),
	}
}

// ============================================================
// 备忘录工具
// ============================================================

// GetNotes 获取所有备忘录
func (d *DailyTools) GetNotes() ([]Note, error) {
	rows, err := d.db.DB.Query(
		"SELECT id, title, content, color, pinned, created_at, updated_at FROM notes ORDER BY pinned DESC, updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		var pinnedInt int
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Color, &pinnedInt, &n.CreatedAt, &n.UpdatedAt); err != nil {
			continue
		}
		n.Pinned = pinnedInt == 1
		notes = append(notes, n)
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历备忘录结果集失败: %w", err)
	}

	return notes, nil
}

// SaveNote 保存备忘录（新增，支持同时设置 pinned 状态）
func (d *DailyTools) SaveNote(title, content, color string, pinned bool) (int64, error) {
	pinnedInt := 0
	if pinned {
		pinnedInt = 1
	}
	result, err := d.db.DB.Exec(
		"INSERT INTO notes (title, content, color, pinned) VALUES (?, ?, ?, ?)",
		title, content, color, pinnedInt,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateNote 更新备忘录内容
func (d *DailyTools) UpdateNote(id int64, title, content, color string) error {
	_, err := d.db.DB.Exec(
		"UPDATE notes SET title=?, content=?, color=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		title, content, color, id,
	)
	return err
}

// PinNote 切换备忘录置顶状态
func (d *DailyTools) PinNote(id int64) error {
	_, err := d.db.DB.Exec(
		"UPDATE notes SET pinned = CASE WHEN pinned = 1 THEN 0 ELSE 1 END WHERE id = ?",
		id,
	)
	return err
}

// DeleteNote 删除备忘录
func (d *DailyTools) DeleteNote(id int64) error {
	_, err := d.db.DB.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}

// ============================================================
// 密码生成工具
// ============================================================

// PasswordResult 密码生成结果
type PasswordResult struct {
	Success  bool   `json:"success"`  // 是否成功
	Password string `json:"password"` // 生成的密码
	Strength string `json:"strength"` // 密码强度（弱/中/强/极强）
	Error    string `json:"error"`    // 错误信息
}

// GeneratePassword 生成随机密码
// length: 密码长度 4-128
// useUpper: 是否包含大写字母
// useLower: 是否包含小写字母
// useNumbers: 是否包含数字
// useSymbols: 是否包含特殊符号
// excludeChars: 需要排除的字符
func (d *DailyTools) GeneratePassword(length int, useUpper, useLower, useNumbers, useSymbols bool, excludeChars string) PasswordResult {
	result := PasswordResult{Success: true}

	// 校验密码长度
	if length < 4 {
		length = 4
	}
	if length > 128 {
		length = 128
	}

	// 至少选择一种字符类型
	if !useUpper && !useLower && !useNumbers && !useSymbols {
		return PasswordResult{Success: false, Error: "至少需要选择一种字符类型"}
	}

	// 构建字符池
	var pool string
	var requiredChars []string // 每种选中类型至少一个字符

	if useUpper {
		upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		pool += upperChars
		// 随机选一个大写字母作为必需字符
		idx := cryptoRandInt(len(upperChars))
		requiredChars = append(requiredChars, string(upperChars[idx]))
	}
	if useLower {
		lowerChars := "abcdefghijklmnopqrstuvwxyz"
		pool += lowerChars
		idx := cryptoRandInt(len(lowerChars))
		requiredChars = append(requiredChars, string(lowerChars[idx]))
	}
	if useNumbers {
		numberChars := "0123456789"
		pool += numberChars
		idx := cryptoRandInt(len(numberChars))
		requiredChars = append(requiredChars, string(numberChars[idx]))
	}
	if useSymbols {
		symbolChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
		pool += symbolChars
		idx := cryptoRandInt(len(symbolChars))
		requiredChars = append(requiredChars, string(symbolChars[idx]))
	}

	// 从字符池中移除排除字符
	if excludeChars != "" {
		filtered := strings.Builder{}
		for _, ch := range pool {
			if !strings.ContainsRune(excludeChars, ch) {
				filtered.WriteRune(ch)
			}
		}
		pool = filtered.String()
		if pool == "" {
			return PasswordResult{Success: false, Error: "排除字符后可用字符池为空"}
		}
		// 重新生成必需字符（确保不在排除列表中）
		var newRequired []string
		for _, rc := range requiredChars {
			if !strings.ContainsAny(rc, excludeChars) {
				newRequired = append(newRequired, rc)
			} else {
				// 从剩余字符池中重新选取
				if len(pool) > 0 {
					idx := cryptoRandInt(len(pool))
					newRequired = append(newRequired, string(pool[idx]))
				}
			}
		}
		requiredChars = newRequired
	}

	// 生成密码
	password := make([]byte, 0, length)

	// 先添加必需字符（确保每种类型至少一个）
	password = append(password, []byte(strings.Join(requiredChars, ""))...)

	// 填充剩余长度
	remaining := length - len(password)
	for i := 0; i < remaining; i++ {
		idx := cryptoRandInt(len(pool))
		password = append(password, pool[idx])
	}

	// 打乱密码字符顺序（Fisher-Yates 洗牌算法）
	for i := len(password) - 1; i > 0; i-- {
		j := cryptoRandInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	result.Password = string(password)

	// 评估密码强度
	result.Strength = evaluatePasswordStrength(result.Password)

	return result
}

// evaluatePasswordStrength 评估密码强度
func evaluatePasswordStrength(password string) string {
	length := len(password)
	score := 0

	// 长度评分
	if length >= 8 {
		score++
	}
	if length >= 12 {
		score++
	}
	if length >= 16 {
		score++
	}

	// 字符类型评分
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSymbol := false
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		default:
			hasSymbol = true
		}
	}

	types := 0
	if hasUpper {
		types++
	}
	if hasLower {
		types++
	}
	if hasNumber {
		types++
	}
	if hasSymbol {
		types++
	}
	score += types

	// 评定强度
	switch {
	case score >= 6:
		return "极强"
	case score >= 4:
		return "强"
	case score >= 3:
		return "中"
	default:
		return "弱"
	}
}

// cryptoRandInt 生成安全的随机整数 [0, max)
func cryptoRandInt(max int) int {
	if max <= 0 {
		return 0
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		// 降级到非安全随机（不应该发生）
		return 0
	}
	return int(n.Int64())
}

// ============================================================
// 文档转换工具（DocConverter）
// ============================================================

// ConvertDocument 文档格式转换
// input: 输入文本
// fromFormat: 源格式（markdown/html/plain）
// toFormat: 目标格式（markdown/html/plain）
func (d *DailyTools) ConvertDocument(input string, fromFormat, toFormat string) ToolResult {
	fromFormat = strings.ToLower(strings.TrimSpace(fromFormat))
	toFormat = strings.ToLower(strings.TrimSpace(toFormat))

	if input == "" {
		return ToolResult{Success: false, Error: "输入内容不能为空"}
	}

	// 校验格式
	validFormats := map[string]bool{"markdown": true, "html": true, "plain": true}
	if !validFormats[fromFormat] || !validFormats[toFormat] {
		return ToolResult{Success: false, Error: "不支持的格式（支持: markdown/html/plain）"}
	}

	// 相同格式直接返回
	if fromFormat == toFormat {
		return ToolResult{Success: true, Data: input}
	}

	// 根据源格式和目标格式选择转换函数
	key := fromFormat + "->" + toFormat
	switch key {
	case "markdown->html":
		return ToolResult{Success: true, Data: convertMarkdownToHTML(input)}
	case "html->plain":
		return ToolResult{Success: true, Data: convertHTMLToPlain(input)}
	case "html->markdown":
		return ToolResult{Success: true, Data: convertHTMLToMarkdown(input)}
	case "plain->html":
		return ToolResult{Success: true, Data: convertPlainToHTML(input)}
	case "plain->markdown":
		return ToolResult{Success: true, Data: convertPlainToMarkdown(input)}
	case "markdown->plain":
		// Markdown -> HTML -> Plain
		htmlContent := convertMarkdownToHTML(input)
		return ToolResult{Success: true, Data: convertHTMLToPlain(htmlContent)}
	default:
		return ToolResult{Success: false, Error: fmt.Sprintf("不支持的转换：%s -> %s", fromFormat, toFormat)}
	}
}

// convertMarkdownToHTML 将 Markdown 转换为 HTML
// 支持：标题(#, ##, ###)、粗体(**text**)、斜体(*text*)、链接[text](url)、
// 代码块(```)、行内代码(`code`)、无序列表(- item)、有序列表(1. item)、分隔线(---)
func convertMarkdownToHTML(md string) string {
	var sb strings.Builder
	lines := strings.Split(md, "\n")
	i := 0

	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// 空行
		if trimmed == "" {
			i++
			continue
		}

		// 代码块（```开头）
		if strings.HasPrefix(trimmed, "```") {
			language := strings.TrimPrefix(trimmed, "```")
			language = strings.TrimSpace(language)
			sb.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">", language))
			i++
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
				sb.WriteString(html.EscapeString(lines[i]))
				sb.WriteString("\n")
				i++
			}
			sb.WriteString("</code></pre>\n")
			if i < len(lines) {
				i++ // 跳过结束的 ```
			}
			continue
		}

		// 分隔线（--- 或 *** 或 ___）
		if trimmed == "---" || trimmed == "***" || trimmed == "___" {
			sb.WriteString("<hr />\n")
			i++
			continue
		}

		// 标题（# ~ ######）
		if strings.HasPrefix(trimmed, "#") {
			level := 0
			for _, ch := range trimmed {
				if ch == '#' {
					level++
				} else {
					break
				}
			}
			if level <= 6 && (len(trimmed) == level || trimmed[level] == ' ') {
				title := strings.TrimSpace(trimmed[level:])
				title = inlineMarkdownToHTML(title)
				sb.WriteString(fmt.Sprintf("<h%d>%s</h%d>\n", level, title, level))
				i++
				continue
			}
		}

		// 无序列表（- 或 * 开头）
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			sb.WriteString("<ul>\n")
			for i < len(lines) {
				itemLine := strings.TrimSpace(lines[i])
				if !strings.HasPrefix(itemLine, "- ") && !strings.HasPrefix(itemLine, "* ") {
					break
				}
				item := strings.TrimPrefix(itemLine, "- ")
				item = strings.TrimPrefix(item, "* ")
				item = inlineMarkdownToHTML(item)
				sb.WriteString(fmt.Sprintf("<li>%s</li>\n", item))
				i++
			}
			sb.WriteString("</ul>\n")
			continue
		}

		// 有序列表（1. 2. 3. 开头）
		if isOrderedListLine(trimmed) {
			sb.WriteString("<ol>\n")
			for i < len(lines) {
				itemLine := strings.TrimSpace(lines[i])
				if !isOrderedListLine(itemLine) {
					break
				}
				// 去掉数字和点
				dotIdx := strings.Index(itemLine, ".")
				if dotIdx < 0 {
					break
				}
				item := strings.TrimSpace(itemLine[dotIdx+1:])
				item = inlineMarkdownToHTML(item)
				sb.WriteString(fmt.Sprintf("<li>%s</li>\n", item))
				i++
			}
			sb.WriteString("</ol>\n")
			continue
		}

		// 普通段落
		para := inlineMarkdownToHTML(trimmed)
		sb.WriteString(fmt.Sprintf("<p>%s</p>\n", para))
		i++
	}

	return sb.String()
}

// isOrderedListLine 判断是否为有序列表行
func isOrderedListLine(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) < 3 {
		return false
	}
	// 必须以数字开头，后跟点和空格
	for i := 0; i < len(line); i++ {
		if line[i] >= '0' && line[i] <= '9' {
			continue
		}
		if line[i] == '.' && i+1 < len(line) && line[i+1] == ' ' {
			return true
		}
		return false
	}
	return false
}

// inlineMarkdownToHTML 处理行内 Markdown 元素
// 支持：粗体(**text**)、斜体(*text*)、行内代码(`code`)、链接[text](url)、图片![alt](url)
func inlineMarkdownToHTML(text string) string {
	// 处理行内代码 `code`
	codeRe := regexp.MustCompile("`([^`]+)`")
	text = codeRe.ReplaceAllString(text, "<code>$1</code>")

	// 处理图片 ![alt](url)
	imgRe := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	text = imgRe.ReplaceAllString(text, `<img src="$2" alt="$1" />`)

	// 处理链接 [text](url)
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRe.ReplaceAllString(text, `<a href="$2">$1</a>`)

	// 处理粗体 **text** 或 __text__
	boldRe := regexp.MustCompile(`\*\*(.+?)\*\*|__(.+?)__`)
	text = boldRe.ReplaceAllStringFunc(text, func(match string) string {
		inner := strings.TrimSuffix(strings.TrimPrefix(match, "**"), "**")
		inner = strings.TrimSuffix(strings.TrimPrefix(inner, "__"), "__")
		return fmt.Sprintf("<strong>%s</strong>", inner)
	})

	// 处理斜体 *text* 或 _text_（注意不匹配粗体）
	italicRe := regexp.MustCompile(`\*(.+?)\*|_(.+?)_`)
	text = italicRe.ReplaceAllStringFunc(text, func(match string) string {
		// 跳过已经是 HTML 标签的内容
		if strings.Contains(match, "<") {
			return match
		}
		inner := strings.TrimSuffix(strings.TrimPrefix(match, "*"), "*")
		inner = strings.TrimSuffix(strings.TrimPrefix(inner, "_"), "_")
		if inner == "" {
			return match
		}
		return fmt.Sprintf("<em>%s</em>", inner)
	})

	return text
}

// convertHTMLToPlain 将 HTML 转换为纯文本
// 使用正则去除所有 HTML 标签，保留文本内容
func convertHTMLToPlain(htmlContent string) string {
	// 去除 script 和 style 标签及其内容
	scriptRe := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	htmlContent = scriptRe.ReplaceAllString(htmlContent, "")

	styleRe := regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	htmlContent = styleRe.ReplaceAllString(htmlContent, "")

	// 将 <br>、<br/>、<p> 等块级标签替换为换行
	blockRe := regexp.MustCompile(`(?i)<br\s*/?>|</p>|</div>|</h[1-6]>|</li>|</tr>|<hr\s*/?>`)
	htmlContent = blockRe.ReplaceAllString(htmlContent, "\n")

	// 将块级开始标签也替换为换行
	blockStartRe := regexp.MustCompile(`(?i)<p[^>]*>|<div[^>]*>|<h[1-6][^>]*>|<li[^>]*>`)
	htmlContent = blockStartRe.ReplaceAllString(htmlContent, "\n")

	// 将列表项标记
	htmlContent = regexp.MustCompile(`(?i)</ul>|</ol>`).ReplaceAllString(htmlContent, "\n")

	// 解码 HTML 实体
	htmlContent = html.UnescapeString(htmlContent)

	// 去除所有剩余 HTML 标签
	tagRe := regexp.MustCompile(`<[^>]+>`)
	htmlContent = tagRe.ReplaceAllString(htmlContent, "")

	// 将多个连续换行合并为最多两个
	newlineRe := regexp.MustCompile(`\n{3,}`)
	htmlContent = newlineRe.ReplaceAllString(htmlContent, "\n\n")

	// 去除行首行尾空白
	lines := strings.Split(htmlContent, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return strings.Join(result, "\n")
}

// convertHTMLToMarkdown 将 HTML 转换为 Markdown（简单逆向转换）
func convertHTMLToMarkdown(htmlContent string) string {
	// 先转为纯文本作为基础
	plain := convertHTMLToPlain(htmlContent)

	// 简单的逆向转换
	// 由于完整的 HTML -> Markdown 转换非常复杂，这里提供基本支持
	// 提取标题
	titleRe := regexp.MustCompile(`(?i)<h([1-6])[^>]*>(.*?)</h[1-6]>`)
	htmlContent = titleRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := titleRe.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		level, _ := strconv.Atoi(sub[1])
		title := convertHTMLToPlain(sub[2])
		return strings.Repeat("#", level) + " " + title + "\n\n"
	})

	// 提取链接
	linkRe := regexp.MustCompile(`(?i)<a[^>]*href="([^"]*)"[^>]*>(.*?)</a>`)
	htmlContent = linkRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := linkRe.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		text := convertHTMLToPlain(sub[2])
		return fmt.Sprintf("[%s](%s)", text, sub[1])
	})

	// 提取粗体
	boldRe := regexp.MustCompile(`(?i)<strong[^>]*>(.*?)</strong>|<b[^>]*>(.*?)</b>`)
	htmlContent = boldRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := boldRe.FindStringSubmatch(match)
		for _, s := range sub[1:] {
			if s != "" {
				return fmt.Sprintf("**%s**", convertHTMLToPlain(s))
			}
		}
		return match
	})

	// 提取斜体
	italicRe := regexp.MustCompile(`(?i)<em[^>]*>(.*?)</em>|<i[^>]*>(.*?)</i>`)
	htmlContent = italicRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := italicRe.FindStringSubmatch(match)
		for _, s := range sub[1:] {
			if s != "" {
				return fmt.Sprintf("*%s*", convertHTMLToPlain(s))
			}
		}
		return match
	})

	// 提取行内代码
	codeRe := regexp.MustCompile(`(?i)<code[^>]*>(.*?)</code>`)
	htmlContent = codeRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := codeRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		return fmt.Sprintf("`%s`", convertHTMLToPlain(sub[1]))
	})

	// 提取图片
	imgRe := regexp.MustCompile(`(?i)<img[^>]*src="([^"]*)"[^>]*alt="([^"]*)"[^>]*/?>|<img[^>]*alt="([^"]*)"[^>]*src="([^"]*)"[^>]*/?>`)
	htmlContent = imgRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := imgRe.FindStringSubmatch(match)
		src, alt := "", ""
		if len(sub) >= 3 {
			src = sub[1]
			alt = sub[2]
		}
		if src == "" && len(sub) >= 5 {
			src = sub[4]
			alt = sub[3]
		}
		return fmt.Sprintf("![%s](%s)", alt, src)
	})

	// 提取代码块
	preRe := regexp.MustCompile(`(?is)<pre[^>]*><code[^>]*>(.*?)</code></pre>`)
	htmlContent = preRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := preRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		code := html.UnescapeString(sub[1])
		return fmt.Sprintf("```\n%s\n```\n", code)
	})

	// 去除剩余标签
	tagRe := regexp.MustCompile(`<[^>]+>`)
	htmlContent = tagRe.ReplaceAllString(htmlContent, "")

	// 清理多余空行
	newlineRe := regexp.MustCompile(`\n{3,}`)
	htmlContent = newlineRe.ReplaceAllString(htmlContent, "\n\n")

	return strings.TrimSpace(htmlContent)
}

// convertPlainToHTML 将纯文本转换为 HTML
// 转义 HTML 实体，换行转 <br>，段落用 <p> 包裹
func convertPlainToHTML(text string) string {
	// 转义 HTML 特殊字符
	escaped := html.EscapeString(text)

	// 按段落分割（空行分隔）
	paragraphs := strings.Split(escaped, "\n\n")

	var sb strings.Builder
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		// 将段落内的换行转为 <br>
		para = strings.ReplaceAll(para, "\n", "<br />\n")
		sb.WriteString(fmt.Sprintf("<p>%s</p>\n", para))
	}

	return sb.String()
}

// convertPlainToMarkdown 将纯文本转换为 Markdown
// 主要是保持原样，确保特殊字符不被误解
func convertPlainToMarkdown(text string) string {
	// 纯文本转 Markdown 基本不需要转换
	// 只需确保没有 Markdown 特殊字符引起意外格式化
	// 这里直接返回原文
	return text
}
