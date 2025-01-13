//
//  CycleView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI

struct CycleView: View {
    @State private var isTimeActive = false
    var body: some View {
        TimerView()
//        StartTimerView(isActive: $isTimeActive)
//        StopTimerView(isActive: $isTimeActive)
    }
}

#Preview {
    CycleView()
}
