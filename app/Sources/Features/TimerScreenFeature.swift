//
//  TimerScreenFeature.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/19.
//
import ComposableArchitecture
import SwiftUI

@Reducer
struct TimerScreenFeature{
    struct State: Equatable {
        var timer: TimerFeature.State
        var evalModal: EvalModalFeature.State?
    }
    
    enum Action {
        case timer(TimerFeature.Action)
        case evalModal(EvalModalFeature.Action)
        
        case showEvalModal
        case dismissEvalModal
    }
    
    var body: some ReducerOf<Self> {
        Scope(state: \.timer, action: \.timer){TimerFeature()}
        .ifLet(\.evalModal, action: \.evalModal) {EvalModalFeature()}
        
        Reduce {state, action in
            switch action {
            case .timer(.phaseCompleted):
                if state.timer.phase == .shortBreak || state.timer.phase == .longBreak {
                    state.evalModal = EvalModalFeature.State(
                        score: 0.5,
                        round: state.timer.round
                    )
                }
                return .none
                
            case .evalModal(.submitEval(let score)):
                state.evalModal = nil
                print("投稿されたスコア:\(score)")
                return .send(.timer(.start))
                
            case .dismissEvalModal:
                state.evalModal = nil
                return .none
                
            default:
                return .none
            }
        }
    }
}
