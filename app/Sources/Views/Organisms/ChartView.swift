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
    // 固定データ
    private let data: [ConcentrationData] = [
        .init(date: Date().addingTimeInterval(-6*86400), score: 65, movingAverage: 63, stdDev: 6),
        .init(date: Date().addingTimeInterval(-5*86400), score: 70, movingAverage: 67, stdDev: 5),
        .init(date: Date().addingTimeInterval(-4*86400), score: 62, movingAverage: 65, stdDev: 8),
        .init(date: Date().addingTimeInterval(-3*86400), score: 60, movingAverage: 64, stdDev: 9),
        .init(date: Date().addingTimeInterval(-2*86400), score: 75, movingAverage: 68, stdDev: 6),
        .init(date: Date().addingTimeInterval(-1*86400), score: 82, movingAverage: 73, stdDev: 4),
        .init(date: Date(), score: 88, movingAverage: 78, stdDev: 3),
    ]

    // 標準偏差帯
    var stdDevArea: [some ChartContent] {
        data.map {
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
        data.map {
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
        data.map {
            PointMark(
                x: .value("日付", $0.date),
                y: .value("集中度", $0.score)
            )
            .foregroundStyle(ColorTheme.navy)
        }
    }

    // 移動平均平均
    var movingAverageLine: [some ChartContent] {
        data.map {
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
            .chartYScale(domain: 0...100)
            .chartXAxis {
                AxisMarks(values: data.map { $0.date }) { date in
                    AxisGridLine()
                    AxisTick()
                    AxisValueLabel(format: .dateTime.day(.twoDigits))
                }
            }
            .chartYAxis {
                AxisMarks(values: Array(stride(from: 0, through: 100, by: 20)))
            }

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
                    Text("Average")
                        .font(.caption)
                }

                HStack(spacing: 6) {
                    RoundedRectangle(cornerRadius: 2)
                        .fill(ColorTheme.Gray.opacity(0.4))
                        .frame(width: 24, height: 10)
                    Text("±1 Std. Deviation")
                        .font(.caption)
                }
            }
            .padding(.horizontal)
            .padding(.top, 8)

        }
    }
}

#Preview {
    ChartView()
}
