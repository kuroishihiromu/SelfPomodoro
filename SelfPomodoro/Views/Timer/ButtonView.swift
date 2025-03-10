//
//  ButtonView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI

protocol ButtonViewProtocol {
    var buttonText: String { get }  // ボタンのテキスト
    var buttonColor: Color { get }  // ボタンの色
    var lightColor: Color { get }   // 光沢色
    var shadowColor: Color { get }  // 影の色
    var radius: CGFloat { get }     // ボタンの角丸 
    var action: () -> Void { get }
    
}

extension ButtonViewProtocol {
    var buttonStyle: some View {
        RoundedRectangle(cornerRadius: radius)
            .fill(
                .shadow(.inner(color: lightColor, radius: 6, x: 4, y: 4))
                .shadow(.inner(color: shadowColor, radius: 6, x:-2, y: -2))
            )
            .foregroundColor(buttonColor)
            .shadow(color: buttonColor, radius: 20, y: 10)
    }
    
    var buttonTextView: some View{
        Text(buttonText)
            .font(.system(size: 20, weight: .semibold, design: .default))
            .foregroundColor(.white)
            .padding(.horizontal, 35)
            .padding(.vertical, 20)
    }
}

struct StartView: View, ButtonViewProtocol {
    var buttonText: String = "Start"
    var buttonColor: Color = Color(red: 0.38, green: 0.28, blue: 0.86)
    var lightColor: Color = Color(red: 0.54, green: 0.41, blue: 0.95)
    var shadowColor: Color = Color(red: 0.25, green: 0.17, blue: 0.75)
    var radius: CGFloat = 25
    var action: () -> Void

    var body: some View {
        Button(action: action) {
            buttonTextView
                .background(buttonStyle)
        }
    }
}

struct LetsTaskView: View, ButtonViewProtocol {
    var buttonText: String = "Let's Task"
    var buttonColor: Color = Color(red: 0.38, green: 0.28, blue: 0.86)
    var lightColor: Color = Color(red: 0.54, green: 0.41, blue: 0.95)
    var shadowColor: Color = Color(red: 0.25, green: 0.17, blue: 0.75)
    var radius: CGFloat = 25
    var action: () -> Void

    var body: some View {
        Button(action: action) {
            buttonTextView
                .background(buttonStyle)
        }
    }
}

struct StartTimerView: View, ButtonViewProtocol {
    var buttonText = "Start Timer"
    var buttonColor = Color(red: 0.38, green: 0.28, blue: 0.86)
    var lightColor = Color(red: 0.54, green: 0.41, blue: 0.95)
    var shadowColor = Color(red: 0.25, green: 0.17, blue: 0.75)
    var radius: CGFloat = 25
    var action: () -> Void
    
    var body: some View {
        Button(action: action) {
            buttonTextView
                .background(buttonStyle)
        }
    }
}

struct StopTimerView: View, ButtonViewProtocol {
    var buttonText = "Stop Timer."
    var buttonColor = Color(red: 0.86, green: 0.28, blue: 0.38)
    var lightColor = Color(red: 0.95, green: 0.41, blue: 0.54)
    var shadowColor = Color(red: 0.75, green: 0.17, blue: 0.25)
    var radius: CGFloat = 25
    var action: () -> Void
    
    var body: some View {
        Button(action: action) {
            buttonTextView
                .background(buttonStyle)
        }
    }
}

struct SetTimerView: View, ButtonViewProtocol {
    var buttonText = "Set Timer."
    var buttonColor = Color(red: 0.38, green: 0.28, blue: 0.86)
    var lightColor = Color(red: 0.54, green: 0.41, blue: 0.95)
    var shadowColor = Color(red: 0.25, green: 0.17, blue: 0.75)
    var radius: CGFloat = 25
    var action: () -> Void
    
    var body: some View {
        Button(action: action) {
            buttonTextView
                .background(buttonStyle)
        }
    }
}
//#Preview {
//    StartTimerView()
//    StopTimerView()
//    SetTimerView()
//}
