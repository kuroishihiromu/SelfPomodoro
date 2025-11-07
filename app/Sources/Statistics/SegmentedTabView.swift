//
//  SegmentedTabView.swift
//  SelfPomodoro
//
//  Created by し on 2025/08/17.
//

import SwiftUI

struct SegmentedTabView: View {
    @Binding var selection: StatisticsFeature.DisplayMode

    var body: some View {
        HStack(spacing: 0) {
            Button { selection = .chart } label: {
                Text("チャート")
                    .font(.headline)
                    .frame(maxWidth: .infinity, minHeight: 40)
                    .foregroundColor(selection == .chart ? ColorTheme.white : ColorTheme.white.opacity(0.7))
                    .background(selection == .chart ? ColorTheme.navy : ColorTheme.Gray.opacity(0.7))
            }

            Button { selection = .heatMap } label: {
                Text("ヒートマップ")
                    .font(.headline)
                    .frame(maxWidth: .infinity, minHeight: 40)
                    .foregroundColor(selection == .heatMap ? ColorTheme.white : ColorTheme.white.opacity(0.7))
                    .background(selection == .heatMap ? ColorTheme.navy : ColorTheme.Gray.opacity(0.7))
            }

            Button { selection = .summary } label: {
                Text("サマリー")
                    .font(.headline)
                    .frame(maxWidth: .infinity, minHeight: 40)
                    .foregroundColor(selection == .summary ? ColorTheme.white : ColorTheme.white.opacity(0.7))
                    .background(selection == .summary ? ColorTheme.navy : ColorTheme.Gray.opacity(0.7))
            }
        }
        .clipShape(RoundedRectangle(cornerRadius: 10))
    }
}
