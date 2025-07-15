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
                Text("Concentration Score Heat Map")
                    .font(.headline)
                    .padding(.horizontal)

                HStack {
                    Button(action: {
                        withAnimation(.easeInOut(duration: 0.3)) {
                            _ = viewStore.send(.previousMonth)
                        }
                    }) {
                        Image(systemName: "chevron.left")
                        Text("Prev")
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
                        Text("Next")
                        Image(systemName: "chevron.right")
                    }
                    .foregroundColor(ColorTheme.black)
                }
                .padding(.horizontal)

                HourlyHeatMapGridView(
                    focusData: viewStore.focusData,
                    currentMonth: viewStore.currentMonth
                )

                HStack {
                    Text("縦軸：日　横軸：時間")
                        .font(.caption2)

                    Spacer()

                    HStack(spacing: 4) {
                        Text("Low")
                            .font(.caption2)

                        ForEach(0..<5, id: \.self) { i in
                            Rectangle()
                                .fill(ColorTheme.navy.opacity(0.2 + Double(i) * 0.2))
                                .frame(width: 16, height: 16)
                                .cornerRadius(2)
                        }

                        Text("High")
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

struct CalendarGridView: View {
    let data: [FocusData]
    let currentMonth: Date

    var body: some View {
        let gridData = HeatMapDataProcessor.generateCalendarData(for: currentMonth)
         VStack(spacing: 4) {
             HStack(spacing: 4) {
                ForEach(0..<7, id: \.self) { index in
                    Text(HeatMapDateFormatter.weekdaySymbols[index])
                        .font(.caption)
                        .frame(width: 32, height: 16)
                        .multilineTextAlignment(.center)
                }
            }

             ForEach(gridData, id: \.self) { week in
                 HStack(spacing: 4) {
                     ForEach(week, id: \.self) { date in
                         if let date {
                             let score = score(for: date)
                             ZStack {
                                 Rectangle()
                                     .fill(HeatMapColorMapper.color(for: score))
                                     .frame(width: 32, height: 32)
                                     .cornerRadius(6)

                                 Text("\(HeatMapDataProcessor.dayNumber(for: date))")
                                     .font(.caption2)
                                     .foregroundColor(score.map { $0 > 60 ? .white : .black } ?? .gray)
                             }
                         } else {
                             Rectangle()
                                 .fill(Color.clear)
                                 .frame(width: 32, height: 32)
                         }
                     }
                 }
             }
        }
    }

    private func score(for date: Date) -> Int? {
        return data.first(where: {
            Calendar.current.isDate($0.date, inSameDayAs: date)
        })?.focus_score
    }
}

struct HourlyHeatMapGridView: View {
    let focusData: [FocusData]
    let currentMonth: Date

    var body: some View {
        let gridData = HeatMapDataProcessor.generateHourlyGridData(focusData, for: currentMonth)
        let sortedDates = gridData.keys.sorted()

        ScrollView(.vertical) {
            VStack(alignment: .leading, spacing: 6) {
                // 上部時間ラベル
                HStack(spacing: 3.5) {
                    Text("")
                        .frame(width: 10)
                    ForEach(0..<13, id: \.self) { bucket in
                        Text("\(bucket * 2)")
                            .font(.caption2)
                            .frame(width: 24, alignment: .leading)
                    }
                }

                // 本体
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
