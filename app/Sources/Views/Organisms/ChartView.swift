//
//  ChartView.swift
//  SelfPomodoro
//
//  Created by し on 2025/05/20.
//

import SwiftUI
import Charts

struct ConcentrationData: Identifiable {
    let id = UUID()
    let date: Date
    let score: Double
    let movingAverage: Double
    let stdDev: Double
}

struct ChartView: View {
    // 今日から何日前かの計算
    static func makeDate(daysAgo: Int) -> Date {
        let calendar = Calendar.current
        let today = calendar.startOfDay(for: Date())
        return calendar.date(byAdding: .day, value: -daysAgo, to: today)!
    }

    // 月曜始まりに変換
    static func startOfCurrentWeek() -> Date {
        let calendar = Calendar.current
        let today = calendar.startOfDay(for: Date())
        let weekday = calendar.component(.weekday, from: today)
        let diff = (weekday + 5) % 7
        return calendar.date(byAdding: .day, value: -diff, to: today)!
    }

    // 固定データ
    private let data: [ConcentrationData] = [
        .init(date: makeDate(daysAgo: 20), score: 55, movingAverage: 54, stdDev: 10),
        .init(date: makeDate(daysAgo: 19), score: 60, movingAverage: 56, stdDev: 9),
        .init(date: makeDate(daysAgo: 18), score: 58, movingAverage: 57, stdDev: 8),
        .init(date: makeDate(daysAgo: 17), score: 62, movingAverage: 59, stdDev: 7),
        .init(date: makeDate(daysAgo: 16), score: 65, movingAverage: 61, stdDev: 6),
        .init(date: makeDate(daysAgo: 15), score: 63, movingAverage: 62, stdDev: 6),
        .init(date: makeDate(daysAgo: 14), score: 68, movingAverage: 64, stdDev: 5),

        .init(date: makeDate(daysAgo: 13), score: 70, movingAverage: 66, stdDev: 5),
        .init(date: makeDate(daysAgo: 12), score: 68, movingAverage: 67, stdDev: 4),
        .init(date: makeDate(daysAgo: 11), score: 73, movingAverage: 68, stdDev: 4),
        .init(date: makeDate(daysAgo: 10), score: 69, movingAverage: 69, stdDev: 3),
        .init(date: makeDate(daysAgo: 9), score: 75, movingAverage: 71, stdDev: 3),
        .init(date: makeDate(daysAgo: 8), score: 72, movingAverage: 72, stdDev: 3),
        .init(date: makeDate(daysAgo: 7), score: 77, movingAverage: 73, stdDev: 2),

        .init(date: makeDate(daysAgo: 6), score: 65, movingAverage: 63, stdDev: 6),
        .init(date: makeDate(daysAgo: 5), score: 70, movingAverage: 67, stdDev: 5),
        .init(date: makeDate(daysAgo: 4), score: 62, movingAverage: 65, stdDev: 8),
        .init(date: makeDate(daysAgo: 3), score: 60, movingAverage: 64, stdDev: 9),
        .init(date: makeDate(daysAgo: 2), score: 75, movingAverage: 68, stdDev: 6),
        .init(date: makeDate(daysAgo: 1), score: 82, movingAverage: 73, stdDev: 4),
        .init(date: makeDate(daysAgo: 0), score: 88, movingAverage: 78, stdDev: 3),
    ]

    @State private var currentWeekStart: Date = ChartView.startOfCurrentWeek()   // 今表示している週の月曜の日付を保持

    // 現在の週のデータを抽出
    private var currentWeekData: [ConcentrationData] {
        let weekEnd = Calendar.current.date(byAdding: .day, value: 6, to: currentWeekStart)!
        return data.filter { $0.date >= currentWeekStart && $0.date <= weekEnd }
    }

    // 標準偏差帯
    var stdDevArea: [some ChartContent] {
        currentWeekData.map {
            AreaMark(
                x: .value("日付", $0.date),
                yStart: .value("下限", $0.movingAverage - $0.stdDev),
                yEnd: .value("上限", $0.movingAverage + $0.stdDev)
            )
            .foregroundStyle(ColorTheme.Gray.opacity(0.4))
            .interpolationMethod(.catmullRom)
        }
    }

    // 集中度の線
    var scoreLine: [some ChartContent] {
        currentWeekData
            .filter { $0.score > 0 }
            .map {
            LineMark(
                x: .value("日付", $0.date),
                y: .value("集中度", $0.score),
                series: .value("系列", "Concentration")
            )
            .foregroundStyle(ColorTheme.navy)
            .lineStyle(.init(lineWidth: 3))
        }
    }

    // 集中度の点
    var scorePoints: [some ChartContent] {
        currentWeekData
            .filter { $0.score > 0 }
            .map {
            PointMark(
                x: .value("日付", $0.date),
                y: .value("集中度", $0.score)
            )
            .foregroundStyle(ColorTheme.navy)
        }
    }

    // 移動平均平均
    var movingAverageLine: [some ChartContent] {
        currentWeekData.map {
            LineMark(
                x: .value("日付", $0.date),
                y: .value("移動平均", $0.movingAverage),
                series: .value("系列", "Average")
            )
            .foregroundStyle(ColorTheme.navy)
            .lineStyle(StrokeStyle(lineWidth: 3, dash: [5]))
            .interpolationMethod(.catmullRom)
        }
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Concentration Score Trend")
                .font(.headline)
                .padding(.horizontal)

            Chart {
                ForEach(stdDevArea.indices, id: \.self) { stdDevArea[$0] }
                ForEach(scoreLine.indices, id: \.self) { scoreLine[$0] }
                ForEach(scorePoints.indices, id: \.self) { scorePoints[$0] }
                ForEach(movingAverageLine.indices, id: \.self) { movingAverageLine[$0] }
            }
            .frame(height: 260)
            .padding(.horizontal)
            .chartYScale(domain: 1...100)
            .chartXAxis {
                AxisMarks(values: makeWeekDates()) { date in
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
                        if value.translation.width > 50 {
                            currentWeekStart = previousWeek(from: currentWeekStart)
                        } else if value.translation.width < -50 {
                            currentWeekStart = nextWeek(from: currentWeekStart)
                        }
                    }
            )

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

    // 週を切り替える
    func previousWeek(from date: Date) -> Date {
        Calendar.current.date(byAdding: .day, value: -7, to: date)!
    }

    func nextWeek(from date: Date) -> Date {
        Calendar.current.date(byAdding: .day, value: 7, to: date)!
    }

    // 横軸の日付を作成
    private func makeWeekDates() -> [Date] {
        return (0..<7).compactMap { offset in
            Calendar.current.date(byAdding: .day, value: offset, to: currentWeekStart)
        }
    }
}

#Preview {
    ChartView()
}
