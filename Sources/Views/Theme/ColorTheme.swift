//
//  ColorTheme.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/12.
//

import SwiftUI

struct ColorTheme {
    static let black = Color(hex: "#171a1f")          // 黒ボタン、テキストの色
    static let navy = Color(hex: "#094067")          // 主なボタンやアイコンのカラー
    static let lightNavy = Color(hex: "#565d6d")     // ボタンの非選択時の色
    static let darkGray = Color(hex: "#5f6c7b")      // ちょっとしたメニューバー
    static let Gray = Color(hex: "#bdc1ca")
    static let lightGray = Color(hex: "#f3f4f6")     // テキストフィールドの背景色
    static let skyBlue = Color(hex: "#99cff6")       // 休憩部分や非表示部分
    static let lightSkyBlue = Color(hex: "#f1f8fe")  // 薄いボタン
    static let white = Color(hex: "#fffffe")         // 各ページの背景色、NavyButtonのテキストカラー
}

// HEXカラーをColorに変換する拡張
extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        
        let r, g, b: Double
        if hex.count == 6 {
            r = Double((int >> 16) & 0xFF) / 255.0
            g = Double((int >> 8) & 0xFF) / 255.0
            b = Double(int & 0xFF) / 255.0
        } else {
            r = 1.0
            g = 1.0
            b = 1.0
        }
        
        self.init(red: r, green: g, blue: b)
    }
}
