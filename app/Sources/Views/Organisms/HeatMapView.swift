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

                    Text(Self.monthFormatter.string(from: viewStore.currentMonth))
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


                CalendarGridView(data: viewStore.focusData, currentMonth: viewStore.currentMonth)
            }
            .onAppear {
                viewStore.send(.fetchHeatMap)
            }
        }
    }

    private static let monthFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy年M月"
        return formatter
    }()
}

struct CalendarGridView: View {
    let data: [FocusData]
    let currentMonth: Date

    private var gridData: [[Date?]] {
        generateCalendarData()
    }

    var body: some View {
         VStack(spacing: 4) {
             HStack(spacing: 4) {
                ForEach(0..<7, id: \.self) { index in
                    Text(weekdaySymbols[index])
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
                                     .fill(score.map { color(for: $0) } ?? Color.gray.opacity(0.1))
                                     .frame(width: 32, height: 32)
                                     .cornerRadius(6)

                                 Text("\(dayNumber(for: date))")
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

    private var weekdaySymbols: [String] {
        let formatter = DateFormatter()
        formatter.locale = Locale(identifier: "ja_JP")
        return formatter.shortWeekdaySymbols
    }

    private func generateCalendarData() -> [[Date?]] {
        let calendar = Calendar.current

        guard let range = calendar.range(of: .day, in: .month, for: currentMonth),
              let firstDate = calendar.date(from: calendar.dateComponents([.year, .month], from: currentMonth)) else {
            return []
        }

        var days: [Date?] = []

        // 月初の曜日補正（週の何日目から始まるか）
        let firstWeekday = calendar.dateComponents([.weekday], from: firstDate).weekday ?? 1
        let weekdayIndex = (firstWeekday - calendar.firstWeekday + 7) % 7
        days += Array(repeating: nil, count: weekdayIndex)

        // 日付を配列に追加
        for day in range {
            if let date = calendar.date(byAdding: .day, value: day - 1, to: firstDate) {
                days.append(date)
            }
        }

        // 月末が7の倍数で終わらなければ空マスを追加
        let remainder = days.count % 7
        if remainder != 0 {
            days += Array(repeating: nil, count: 7 - remainder)
        }

        // 週ごとに7日単位で分割
        return stride(from: 0, to: days.count, by: 7).map {
            Array(days[$0..<min($0 + 7, days.count)])
        }
    }

    private func score(for date: Date) -> Int? {
        return data.first(where: {
            Calendar.current.isDate($0.date, inSameDayAs: date)
        })?.focus_score
    }

    private func color(for score: Int) -> Color {
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

    private func dayNumber(for date: Date) -> Int {
        return Calendar.current.component(.day, from: date)
    }
}
