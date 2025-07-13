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


                CalendarGridView(data: viewStore.focusData, currentMonth: viewStore.currentMonth)
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
