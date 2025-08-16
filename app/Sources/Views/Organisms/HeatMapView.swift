//
//  HeatMapView.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/13.
//

import SwiftUI
import ComposableArchitecture

struct HeatMapView: View {
    let store: StoreOf<HeatMapFeature>

    var body: some View {
        WithViewStore(store, observe: { $0 }) { viewStore in
            VStack(spacing: 16) {
                Text(L10n.HeatMap.title)
                    .font(.headline)
                    .padding(.horizontal)

                // 月選択
                HStack {
                    Button(action: {
                        withAnimation(.easeInOut(duration: 0.3)) {
                            _ = viewStore.send(.previousMonth)
                        }
                    }) {
                        Image(systemName: "chevron.left")
                        Text(L10n.HeatMap.prev)
                    }
                    .foregroundColor(ColorTheme.black)

                    Spacer()

                    Text(HeatMapDateFormatter.monthYear.string(from: viewStore.currentMonth))
                        .font(.subheadline)

                    Spacer()

                    Button(action: {
                        withAnimation(.easeInOut(duration: 0.3)) {
                            _ = viewStore.send(.nextMonth)
                        }
                    }) {
                        Text(L10n.HeatMap.next)
                        Image(systemName: "chevron.right")
                    }
                    .foregroundColor(ColorTheme.black)
                }
                .padding(.horizontal)

                HourlyHeatMapGridView(
                    focusData: viewStore.focusData,
                    currentMonth: viewStore.currentMonth
                )

                // 凡例
                HStack {
                    Text(L10n.HeatMap.axisLabel)
                        .font(.caption2)

                    Spacer()

                    HStack(spacing: 4) {
                        Text(L10n.HeatMap.low)
                            .font(.caption2)

                        ForEach(0..<5, id: \.self) { i in
                            Rectangle()
                                .fill(ColorTheme.navy.opacity(0.2 + Double(i) * 0.2))
                                .frame(width: 16, height: 16)
                                .cornerRadius(2)
                        }

                        Text(L10n.HeatMap.high)
                            .font(.caption2)
                    }
                }
                .padding(.horizontal)
                .padding(.bottom, 16)
            }
            .onAppear {
                viewStore.send(.fetchHeatMap)
            }
        }
    }
}

// ヒートマップ本体
struct HourlyHeatMapGridView: View {
    let focusData: [FocusData]
    let currentMonth: Date

    var body: some View {
        let gridData = HeatMapDataProcessor.generateHourlyGridData(focusData, for: currentMonth)
        let sortedDates = gridData.keys.sorted()

        ScrollView(.vertical) {
            VStack(alignment: .leading, spacing: 6) {
                // 横軸
                HStack(spacing: 3.5) {
                    Text("")
                        .frame(width: 10)
                    ForEach(0..<13, id: \.self) { bucket in
                        Text("\(bucket * 2)")
                            .font(.caption2)
                            .frame(width: 24, alignment: .leading)
                    }
                }

                // 縦軸+マス
                ForEach(sortedDates, id: \.self) { date in
                    HStack(spacing: 2) {
                        Text(HeatMapDateFormatter.day.string(from: date))
                            .font(.caption2)
                            .frame(width: 15, alignment: .leading)

                        ForEach(0..<24, id: \.self) { bucket in
                            Rectangle()
                                .fill(HeatMapColorMapper.color(for: gridData[date]?[bucket] ?? nil))
                                .frame(width: 12, height: 12)
                                .cornerRadius(2)
                        }
                    }
                }
            }
        }
        .frame(height: 204)
    }
}

#Preview {
    HeatMapView(
        store: Store(initialState: HeatMapFeature.State()) {
            HeatMapFeature()
        }
    )
}
