//
//  HeatMapper.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/13.
//

import Dependencies
import Foundation

struct HeatMapper {
    var loadFocusData: () -> [FocusData]
}

extension HeatMapper: DependencyKey {
    static let liveValue = HeatMapper {
        guard let url = Bundle.main.url(forResource: "heatmap_data", withExtension: "json"),
              let data = try? Data(contentsOf: url) else {
            return []
        }

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .formatted(HeatMapDateFormatter.yyyyMMdd)

        if let decoded = try? decoder.decode(FocusDataWrapper.self, from: data) {
            return decoded.items
        }

        return []
    }
}

extension DependencyValues {
    var heatMapper: HeatMapper {
        get { self[HeatMapper.self] }
        set { self[HeatMapper.self] = newValue }
    }
}
