//
//  ChartView.swift
//  SelfPomodoro
//
//  Created by し on 2025/05/20.
//

import SwiftUI
import Charts

struct ChartView: View {
    let title: String
    let state: ChartFeature.State
    let sendActionWithAnimation: (ChartFeature.Action) -> Void

    private let swipeThreshold: CGFloat = 50

    private var weekRangeText: String {
        let formatter = DateFormatter()
        formatter.dateFormat = "MM/dd"
        let endDate = Calendar.current.date(byAdding: .day, value: 6, to: state.currentWeekStart)!
        return "\(formatter.string(from: state.currentWeekStart))〜\(formatter.string(from: endDate))"
    }

    var body: some View {
        VStack(spacing: 16) {
            header
            chartBody
            legend
        }
    }

    private var header: some View {
        VStack(spacing: 8) {
            Text(title)
                .font(.headline)
                .padding(.horizontal)

            HStack {
                Button {
                    withAnimation(.easeInOut) {
                        sendActionWithAnimation(.previousWeek)
                    }
                } label: {
                    Image(systemName: "chevron.left")
                    Text("Prev")
                }
                .foregroundColor(ColorTheme.black)

                Spacer()

                Text(weekRangeText)
                    .font(.subheadline)
                    .foregroundColor(ColorTheme.black)

                Spacer()

                Button {
                    withAnimation(.easeInOut) {
                        sendActionWithAnimation(.nextWeek)
                    }
                } label: {
                    Text("Next")
                    Image(systemName: "chevron.right")
                }
                .foregroundColor(ColorTheme.black)
            }
            .padding(.horizontal)
        }
    }

    private var chartBody: some View {
        Chart {
            ForEach(state.weekDates, id: \.self) { date in
                LineMark(
                    x: .value("日付", date),
                    y: .value("透明", 0)
                )
                .foregroundStyle(.clear)
            }

            ForEach(state.currentWeekData) { data in
                if data.score > 0 {
                    AreaMark(
                        x: .value("日付", data.date),
                        yStart: .value("下限", data.movingAverage - data.stdDev),
                        yEnd: .value("上限", data.movingAverage + data.stdDev)
                    )
                    .foregroundStyle(ColorTheme.Gray.opacity(0.4))
                    .interpolationMethod(.catmullRom)
                }
            }

            ForEach(state.currentWeekData) { data in
                if data.score > 0 {
                    LineMark(
                        x: .value("日付", data.date),
                        y: .value("集中度", data.score),
                        series: .value("系列", "Concentration")
                    )
                    .foregroundStyle(ColorTheme.navy)
                    .lineStyle(.init(lineWidth: 3))

                    PointMark(
                        x: .value("日付", data.date),
                        y: .value("集中度", data.score)
                    )
                    .foregroundStyle(ColorTheme.navy)

                    LineMark(
                        x: .value("日付", data.date),
                        y: .value("移動平均", data.movingAverage),
                        series: .value("系列", "Average")
                    )
                    .foregroundStyle(ColorTheme.navy)
                    .lineStyle(StrokeStyle(lineWidth: 3, dash: [5]))
                    .interpolationMethod(.catmullRom)
                }
            }
        }
        .frame(height: 260)
        .chartYScale(domain: 1...100)
        .chartXAxis {
            AxisMarks(values: state.weekDates) { date in
                AxisGridLine()
                AxisTick()
                AxisValueLabel(format: .dateTime.month(.twoDigits).day(.twoDigits))
            }
        }
        .chartYAxis {
            AxisMarks(values: Array(stride(from: 0, through: 100, by: 20)))
        }
        .gesture(
            DragGesture()
                .onEnded { value in
                    withAnimation(.easeInOut) {
                        if value.translation.width > swipeThreshold {
                            sendActionWithAnimation(.previousWeek)
                        } else if value.translation.width < -swipeThreshold {
                            sendActionWithAnimation(.nextWeek)
                        }
                    }
                }
        )
        .padding(.horizontal)
    }

    private var legend: some View {
        HStack(spacing: 16) {
            HStack(spacing: 6) {
                RoundedRectangle(cornerRadius: 2)
                    .fill(ColorTheme.navy)
                    .frame(width: 24, height: 4)
                Text("Concentration")
                    .font(.caption)
            }

            HStack(spacing: 6) {
                RoundedRectangle(cornerRadius: 2)
                    .stroke(ColorTheme.navy, style: StrokeStyle(lineWidth: 2, dash: [5]))
                    .frame(width: 24, height: 4)
                Text("MovingAverage")
                    .font(.caption)
                    .minimumScaleFactor(0.7)
                    .lineLimit(1)
            }

            HStack(spacing: 6) {
                RoundedRectangle(cornerRadius: 2)
                    .fill(ColorTheme.Gray.opacity(0.4))
                    .frame(width: 24, height: 10)
                Text("±1 Std. Deviation")
                    .font(.caption)
                    .minimumScaleFactor(0.7)
                    .lineLimit(1)
            }
        }
        .padding(.horizontal)
    }
}
