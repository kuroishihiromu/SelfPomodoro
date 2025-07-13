//
//  HeatMapPreviewVIew.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/13.
//

import SwiftUI
import ComposableArchitecture

struct HeatMapPreviewView: View {
    var body: some View {
        HeatMapView(
            store: Store(
                initialState: HeatMapFeature.State(),
                reducer: {
                    HeatMapFeature()
                }
            )
        )
    }
}

#Preview {
    HeatMapPreviewView()
}
