//
//  EvalModalFeature.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/19.
//

import ComposableArchitecture
import SwiftUI

@Reducer
struct EvalModalFeature {
    struct State: Equatable {
        var score: Double
        var round: Int
    }
    
    enum Action {
        case updateScore(Double)
        case submitEval(Double)
        case cancel
    }
    
    func reduce(into state: inout State, action: Action) -> Effect<Action> {
        switch action {
        case .updateScore(let value):
            state.score = value
            return .none
        case .submitEval:
            return .none
        case .cancel:
            return .none
        }
    }
}
