import SwiftUI

struct NormalButton: View {
    let text: String
    let bgColor: Color
    let fontColor: Color
    let width: CGFloat
    let height: CGFloat
    let icon: Image?
    let action: () -> Void

    init(
        text: String,
        bgColor: Color,
        fontColor: Color,
        icon: Image? = nil,
        width: CGFloat,
        height: CGFloat,
        action: @escaping () -> Void
    ) {
        self.text = text
        self.bgColor = bgColor
        self.fontColor = fontColor
        self.icon = icon
        self.width = width
        self.height = height
        self.action = action
    }

    var body: some View {
        Button(action: action) {
            HStack(spacing: 15) {
                if let icon = icon {
                    icon
                        .resizable()
                        .frame(width: 18, height: 18)
                }

                Text(text)
                    .foregroundColor(fontColor)
                    .font(.headline)
            }
            .padding()
            .frame(width: width, height: height)
            .background(bgColor)
            .cornerRadius(10)
        }
    }
}

#Preview {
    VStack {
        NormalButton(
            text: "Get Started",
            bgColor: ColorTheme.navy,
            fontColor: ColorTheme.white,
            width: 350,
            height: 52,
            action: { print("yahhoi!") }
        )
        NormalButton(
            text: "Sign Up ",
            bgColor: ColorTheme.navy,
            fontColor: ColorTheme.white,
            width: 350,
            height: 44,
            action: { print("yahhoi!") }
        )
        NormalButton(
            text: "Save",
            bgColor: ColorTheme.navy,
            fontColor: ColorTheme.white,
            width: 160,
            height: 44,
            action: { print("yahhoi!") }
        )
        NormalButton(
            text: "Reset",
            bgColor: ColorTheme.lightSkyBlue,
            fontColor: ColorTheme.navy,
            width: 160,
            height: 44,
            action: { print("yahhoi!") }
        )
        NormalButton(
            text: "Start",
            bgColor: ColorTheme.black,
            fontColor: ColorTheme.white,
            width: 160,
            height: 44,
            action: { print("yahhoi!") }
        )
        NormalButton(
            text: "Add Task",
            bgColor: ColorTheme.navy,
            fontColor: ColorTheme.white,
            icon: Image(.add),
            width: 160,
            height: 44,
            action: { print("yahhoi!") }
        )
    }
}
