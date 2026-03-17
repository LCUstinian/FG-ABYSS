package services

import (
	"crypto/rand"
	"fg-abyss/internal/domain/entity"
	"fmt"
	"math/big"
	"strings"
)

// CodeObfuscator 代码混淆器
type CodeObfuscator struct {
	level entity.ObfuscationLevel
}

// NewCodeObfuscator 创建混淆器
func NewCodeObfuscator(level entity.ObfuscationLevel) *CodeObfuscator {
	return &CodeObfuscator{
		level: level,
	}
}

// Obfuscate 混淆代码
func (o *CodeObfuscator) Obfuscate(code string, language string) (string, error) {
	switch o.level {
	case entity.ObfuscationNone:
		return code, nil
	case entity.ObfuscationLow:
		return o.obfuscateLow(code, language)
	case entity.ObfuscationMedium:
		return o.obfuscateMedium(code, language)
	case entity.ObfuscationHigh:
		return o.obfuscateHigh(code, language)
	default:
		return code, nil
	}
}

// obfuscateLow 轻度混淆
func (o *CodeObfuscator) obfuscateLow(code string, language string) (string, error) {
	result := code
	
	// 变量名混淆
	result = o.obfuscateVariables(result, language)
	
	// 添加垃圾代码
	result = o.addGarbageCode(result, language)
	
	return result, nil
}

// obfuscateMedium 中度混淆
func (o *CodeObfuscator) obfuscateMedium(code string, language string) (string, error) {
	result := code
	
	// 变量名混淆
	result = o.obfuscateVariables(result, language)
	
	// 字符串编码
	result = o.encodeStrings(result, language)
	
	// 添加垃圾代码
	result = o.addGarbageCode(result, language)
	
	// 控制流平坦化（简化版）
	result = o.flattenControlFlow(result, language)
	
	return result, nil
}

// obfuscateHigh 高度混淆
func (o *CodeObfuscator) obfuscateHigh(code string, language string) (string, error) {
	result := code
	
	// 变量名混淆
	result = o.obfuscateVariables(result, language)
	
	// 字符串编码
	result = o.encodeStrings(result, language)
	
	// 添加多层垃圾代码
	result = o.addGarbageCode(result, language)
	result = o.addGarbageCode(result, language)
	
	// 控制流平坦化
	result = o.flattenControlFlow(result, language)
	
	// 代码压缩
	result = o.compressCode(result, language)
	
	return result, nil
}

// obfuscateVariables 混淆变量名
func (o *CodeObfuscator) obfuscateVariables(code string, language string) string {
	// 简化实现：替换常见变量名
	replacements := map[string]string{
		"command":  o.randomVarName(),
		"output":   o.randomVarName(),
		"path":     o.randomVarName(),
		"file":     o.randomVarName(),
		"data":     o.randomVarName(),
		"result":   o.randomVarName(),
		"content":  o.randomVarName(),
		"password": o.randomVarName(),
	}
	
	result := code
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	
	return result
}

// encodeStrings 编码字符串
func (o *CodeObfuscator) encodeStrings(code string, language string) string {
	switch language {
	case "php":
		return o.encodeStringsPHP(code)
	case "asp", "aspx":
		return code // 简化处理
	case "jsp":
		return code // 简化处理
	default:
		return code
	}
}

// encodeStringsPHP PHP 字符串编码
func (o *CodeObfuscator) encodeStringsPHP(code string) string {
	// 查找字符串字面量并 base64 编码
	// 简化实现：只处理简单的字符串
	return code
}

// addGarbageCode 添加垃圾代码
func (o *CodeObfuscator) addGarbageCode(code string, language string) string {
	garbage := o.generateGarbageCode(language)
	
	// 在代码开头插入垃圾代码
	return garbage + "\n" + code
}

// generateGarbageCode 生成垃圾代码
func (o *CodeObfuscator) generateGarbageCode(language string) string {
	switch language {
	case "php":
		return o.garbagePHP()
	case "asp", "aspx":
		return o.garbageASP()
	case "jsp":
		return o.garbageJSP()
	default:
		return ""
	}
}

// garbagePHP PHP 垃圾代码
func (o *CodeObfuscator) garbagePHP() string {
	return fmt.Sprintf(`<?php
// Obfuscation Layer %d
$_%s = %d;
$_%s = "%s";
?>`,
		o.randomInt(1000, 9999),
		o.randomVarName(), o.randomInt(0, 100),
		o.randomVarName(), o.randomString(10),
	)
}

// garbageASP ASP 垃圾代码
func (o *CodeObfuscator) garbageASP() string {
	return fmt.Sprintf(`<?
' Obfuscation Layer %d
Dim _%s
_%s = %d
?>`,
		o.randomInt(1000, 9999),
		o.randomVarName(),
		o.randomVarName(),
		o.randomInt(0, 100),
	)
}

// garbageJSP JSP 垃圾代码
func (o *CodeObfuscator) garbageJSP() string {
	return fmt.Sprintf(`<?
// Obfuscation Layer %d
String _%s = "%s";
int _%s = %d;
?>`,
		o.randomInt(1000, 9999),
		o.randomVarName(), o.randomString(10),
		o.randomVarName(), o.randomInt(0, 100),
	)
}

// flattenControlFlow 控制流平坦化（简化版）
func (o *CodeObfuscator) flattenControlFlow(code string, language string) string {
	// 简化实现：添加额外的条件判断
	switch language {
	case "php":
		return fmt.Sprintf(`<?php
if (true) {
%s
}
?>`, code)
	default:
		return code
	}
}

// compressCode 代码压缩
func (o *CodeObfuscator) compressCode(code string, language string) string {
	// 移除空白和注释
	lines := strings.Split(code, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "#") {
			result = append(result, trimmed)
		}
	}
	
	return strings.Join(result, "\n")
}

// randomVarName 生成随机变量名
func (o *CodeObfuscator) randomVarName() string {
	prefix := "_"
	chars := "abcdefghijklmnopqrstuvwxyz"
	length := 8
	
	result := prefix
	for i := 0; i < length; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result += string(chars[idx.Int64()])
	}
	
	return result
}

// randomString 生成随机字符串
func (o *CodeObfuscator) randomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	
	for i := 0; i < length; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result += string(chars[idx.Int64()])
	}
	
	return result
}

// randomInt 生成随机整数
func (o *CodeObfuscator) randomInt(min, max int64) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return int(n.Int64() + min)
}

// ObfuscatePHP PHP 代码混淆（便捷方法）
func ObfuscatePHP(code string, level entity.ObfuscationLevel) (string, error) {
	obfuscator := NewCodeObfuscator(level)
	return obfuscator.Obfuscate(code, "php")
}

// ObfuscateASP ASP 代码混淆（便捷方法）
func ObfuscateASP(code string, level entity.ObfuscationLevel) (string, error) {
	obfuscator := NewCodeObfuscator(level)
	return obfuscator.Obfuscate(code, "asp")
}

// ObfuscateJSP JSP 代码混淆（便捷方法）
func ObfuscateJSP(code string, level entity.ObfuscationLevel) (string, error) {
	obfuscator := NewCodeObfuscator(level)
	return obfuscator.Obfuscate(code, "jsp")
}
