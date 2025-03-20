import SwiftUI

protocol ButtonViewProtocol {
    var buttonText: String { get }  // ボタンのテキスト
    var radius: CGFloat { get }     // ボタンの角丸
    var action: () -> Void { get }
}

extension ButtonViewProtocol {
    var buttonStyle: some View {
        RoundedRectangle(cornerRadius: radius)
            .foregroundColor(ColorTheme.black)
    }
    
    var buttonTextView: some View {
        Text(buttonText)
            .font(.system(size: 20, weight: .semibold, design: .default))
            .foregroundColor(.white)
            .padding(.horizontal, 35)
            .padding(.vertical, 20)
    }
}

struct StartView: View, ButtonViewProtocol {
    var buttonText: String = "Start"
    var colorScheme = ColorTheme.navy
    var radius: CGFloat = 10
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
    var colorScheme = ColorTheme.navy
    var radius: CGFloat = 10
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
    var colorScheme = ColorTheme.skyBlue
    var radius: CGFloat = 10
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
    var colorScheme = ColorTheme.navy
    var radius: CGFloat = 10
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
    var colorScheme = ColorTheme.black
    var radius: CGFloat = 10
    var action: () -> Void
    
    var body: some View {
        Button(action: action) {
            buttonTextView
                .background(buttonStyle)
        }
    }
}

#Preview {
    VStack {
        StartView(action: {})
        LetsTaskView(action: {})
        StartTimerView(action: {})
        StopTimerView(action: {})
        SetTimerView(action: {})
    }
}

