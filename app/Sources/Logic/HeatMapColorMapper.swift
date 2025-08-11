//
//  HeatMapColorMapper.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/13.
//

import SwiftUI

enum HeatMapColorMapper {
    static func color(for score: Int?) -> Color {
        guard let score = score else {
            return Color.gray.opacity(0.05)
        }

        switch score {
        case 91...100: return ColorTheme.navy
        case 81...90:  return ColorTheme.navy.opacity(0.9)
        case 71...80:  return ColorTheme.navy.opacity(0.8)
        case 61...70:  return ColorTheme.navy.opacity(0.7)
        case 51...60:  return ColorTheme.navy.opacity(0.6)
        case 41...50:  return ColorTheme.navy.opacity(0.5)
        case 31...40:  return ColorTheme.navy.opacity(0.4)
        case 21...30:  return ColorTheme.navy.opacity(0.3)
        case 11...20:  return ColorTheme.navy.opacity(0.2)
        case 1...10:   return ColorTheme.navy.opacity(0.1)
        default:       return Color.gray.opacity(0.05)
        }
    }
}
