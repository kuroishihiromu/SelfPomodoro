//
//  ChartView.swift
//  SelfPomodoro
//
//  Created by し on 2025/05/20.
//

import SwiftUI
import Charts
import Dependencies

struct ChartView: View {
    @Dependency(\.statisticsAPIClient) var statisticsAPIClient

    @State private var data: [ConcentrationData] = []   // グラフに表示するデータのリスト
    @State private var currentWeekStart: Date = ChartView.startOfCurrentWeek()   // 現在表示している週の月曜の日付を保持

    // グラフの横軸表示用に月曜起点で8日分の日付を生成
    private var weekDates: [Date] {
        (0..<8).compactMap { offset in
            Calendar.current.date(byAdding: .day, value: offset, to: currentWeekStart)
        }
    }

    // 現在の週のデータを抽出
    private var currentWeekData: [ConcentrationData] {
        let dict = Dictionary(uniqueKeysWithValues: data.map { ($0.date, $0) })
        return weekDates.map { date in
            dict[date] ?? ConcentrationData(date: date, score: 0, movingAverage: 0, stdDev: 0)
        }
    }

    // 標準偏差帯
    var stdDevArea: [some ChartContent] {
        currentWeekData
            .filter { $0.score > 0 }
            .map {
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
        currentWeekData
            .filter { $0.score > 0 }
            .map {
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
                // 透明のラインを描画しておくことで、ChartのX軸にweekDates全体を認識させる
                ForEach(weekDates, id: \.self) { date in
                    LineMark(
                        x: .value("日付", date),
                        y: .value("透明", 0)
                    )
                    .foregroundStyle(.clear)
                }
                ForEach(stdDevArea.indices, id: \.self) { stdDevArea[$0] }
                ForEach(scoreLine.indices, id: \.self) { scoreLine[$0] }
                ForEach(scorePoints.indices, id: \.self) { scorePoints[$0] }
                ForEach(movingAverageLine.indices, id: \.self) { movingAverageLine[$0] }
            }
            .frame(height: 260)
            .padding(.horizontal)
            .chartYScale(domain: 1...100)
            .chartXAxis {
                AxisMarks(values: weekDates) { date in
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
                            if value.translation.width > 50 {
                                currentWeekStart = previousWeek(from: currentWeekStart)
                            } else if value.translation.width < -50 {
                                currentWeekStart = nextWeek(from: currentWeekStart)
                            }
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
        .task {
            await loadData()
        }
    }

    // データ取得
    private func loadData() async {
        do {
            let result = try await statisticsAPIClient.fetchFocusTrend()

            let normalized = result.map { r in
                FocusTrendResult(date: Calendar.current.startOfDay(for: r.date), focusScore: r.focusScore)
            }

            self.data = calculateMovingAverage(from: normalized)
            print("取得成功: \(self.data.count) 件")
        } catch {
            print("データ取得失敗: \(error.localizedDescription)")
        }
    }

    // 月曜始まりに変換
    static func startOfCurrentWeek() -> Date {
        let calendar = Calendar.current
        let today = calendar.startOfDay(for: Date())
        let weekday = calendar.component(.weekday, from: today)
        let diff = (weekday + 5) % 7
        return calendar.date(byAdding: .day, value: -diff, to: today)!
    }

    // 週を切り替える
    func previousWeek(from date: Date) -> Date {
        Calendar.current.date(byAdding: .day, value: -7, to: date)!
    }

    func nextWeek(from date: Date) -> Date {
        Calendar.current.date(byAdding: .day, value: 7, to: date)!
    }
}

#Preview {
    ChartView()
}
