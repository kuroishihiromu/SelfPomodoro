import SwiftUI

protocol ButtonViewProtocol {
    var buttonText: String { get }  // ボタンのテキスト
    var colorScheme: AppColorScheme { get }  // カラースキーム
    var radius: CGFloat { get }     // ボタンの角丸
    var action: () -> Void { get }
}

extension ButtonViewProtocol {
    var buttonStyle: some View {
        RoundedRectangle(cornerRadius: radius)
            .fill(
                .shadow(.inner(color: colorScheme.light, radius: 6, x: 4, y: 4))
                .shadow(.inner(color: colorScheme.shadow, radius: 6, x: -2, y: -2))
            )
            .foregroundColor(colorScheme.main)
            .shadow(color: colorScheme.main, radius: 20, y: 10)
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
    var colorScheme = Color.theme.primary  // `primary` カラーを適用
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
    var colorScheme = Color.theme.secondary  // `secondary` カラーを適用
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
    var colorScheme = Color.theme.accent  // `accent` カラーを適用
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
    var colorScheme = Color.theme.highlight  // `highlight` カラーを適用
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
    var colorScheme = Color.theme.primary  // `primary` カラーを適用
    var radius: CGFloat = 25
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

