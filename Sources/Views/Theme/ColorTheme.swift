//
//  ColorTheme.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/12.
//

import SwiftUI

/// `mainColor` を基準に `lightColor` と `shadowColor` を算出する
struct AppColorScheme {
    let main: Color
    let light: Color
    let shadow: Color

    init(main: Color) {
        self.main = main
        self.light = main.lightened(by: 0.2)  // 20% 明るくする
        self.shadow = main.darkened(by: 0.2)  // 20% 暗くする
    }
}

struct ColorTheme {
    /// **Elements**
    static let background = AppColorScheme(main: Color(hex: "#fffffe"))  // 背景
    static let headline = AppColorScheme(main: Color(hex: "#094067"))  // ヘッドライン（タイトル）
    static let paragraph = AppColorScheme(main: Color(hex: "#5f6c7b"))  // 段落（本文）
    static let button = AppColorScheme(main: Color(hex: "#3da9fc"))  // ボタン
    static let buttonText = AppColorScheme(main: Color(hex: "#fffffe"))  // ボタンテキスト

    /// **Illustration**
    static let stroke = AppColorScheme(main: Color(hex: "#094067"))  // 図のアウトライン
    static let illustrationMain = AppColorScheme(main: Color(hex: "#fffffe"))  // イラストのメイン色
    static let highlight = AppColorScheme(main: Color(hex: "#3da9fc"))  // 強調色
    static let secondary = AppColorScheme(main: Color(hex: "#90b4ce"))  // セカンダリーカラー
    static let tertiary = AppColorScheme(main: Color(hex: "#ef4565"))  // 補助的な色（警告やアラート）

}

extension Color {
    static let theme = ColorTheme.self
    
    /// HEX 文字列から Color を生成する
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let r, g, b: Double
        r = Double((int >> 16) & 0xFF) / 255.0
        g = Double((int >> 8) & 0xFF) / 255.0
        b = Double(int & 0xFF) / 255.0
        self.init(red: r, green: g, blue: b)
    }

    /// 既存のカラーを明るくする
    func lightened(by percentage: CGFloat) -> Color {
        return self.adjustBrightness(by: abs(percentage))
    }

    /// 既存のカラーを暗くする
    func darkened(by percentage: CGFloat) -> Color {
        return self.adjustBrightness(by: -abs(percentage))
    }

    /// カラーの明るさを調整する
    private func adjustBrightness(by percentage: CGFloat) -> Color {
        var r: CGFloat = 0, g: CGFloat = 0, b: CGFloat = 0, a: CGFloat = 0
        UIColor(self).getRed(&r, green: &g, blue: &b, alpha: &a)
        return Color(red: min(r + percentage, 1.0),
                     green: min(g + percentage, 1.0),
                     blue: min(b + percentage, 1.0),
                     opacity: Double(a))
    }
}
