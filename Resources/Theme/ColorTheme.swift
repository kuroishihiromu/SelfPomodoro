//
//  ColorTheme.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/12.
//

import SwiftUI

/// `mainColor` を元に `lightColor` と `shadowColor` を算出する
struct AppColorScheme {
    let main: Color
    let light: Color
    let shadow: Color

    init(main: Color) {
        self.main = main
        self.light = main.lightened(by: 0.2)  // 20%明るくする
        self.shadow = main.darkened(by: 0.2)  // 20%暗くする
    }
}

struct ColorTheme {
    /// 主要なカラースキーム（main を基準に light, shadow を算出）
    static let primary = AppColorScheme(main: Color(hex: "#211C84"))  // 背景・大枠
    static let secondary = AppColorScheme(main: Color(hex: "#4D55CC"))  // 説明・コンテナ
    static let accent = AppColorScheme(main: Color(hex: "#7A73D1"))  // アクセントカラー
    static let highlight = AppColorScheme(main: Color(hex: "#B5A8D5"))  // 協調色

    static let background = Color(hex: "#E5E7F5")  // アプリ全体の背景色
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
