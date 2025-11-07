//
//  SessionAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation
import Dependencies

enum SessionAPIError: Error, Equatable {
    case networkError
    case decodingError
    case unknown
}

struct SessionAPIClient {
    var startSession: () async throws -> SessionResult
    var completeSession: (_ sessionId: UUID) async throws -> SessionResult
    var startRound: (_ sessionId: UUID) async throws -> RoundResult
    var completeRound: (_ roundId: UUID, _ focusScore: Int) async throws -> RoundResult
    var getSession: (_ sessionId: UUID) async throws -> SessionResult
}

extension SessionAPIClient {
    private static func mockSession(id: UUID = UUID(), completed: Bool = false) -> SessionResult {
        SessionResult(
            id: id,
            startTime: Date().addingTimeInterval(-600),
            endTime: completed ? Date() : nil,
            averageFocus: completed ? 0.7 : nil,
            totalWorkMin: completed ? 60 : nil,
            roundCount: completed ? 4 : nil,
            breakTime: completed ? 15 : nil
        )
    }

    private static func mockRound(sessionId: UUID, order: Int) -> RoundResult {
        RoundResult(
            id: UUID(),
            sessionId: sessionId,
            roundOrder: order,
            startTime: Date(),
            endTime: nil,
            workTime: nil,
            breakTime: nil,
            focusScore: nil
        )
    }

    private static func mockRoundCompletion(roundId: UUID, focusScore: Int) -> RoundResult {
        RoundResult(
            id: roundId,
            sessionId: UUID(),
            roundOrder: 1,
            startTime: Date().addingTimeInterval(-1500),
            endTime: Date(),
            workTime: 25,
            breakTime: 5,
            focusScore: focusScore
        )
    }

    static let live = SessionAPIClient(
        startSession: {
            SessionAPIClient.mockSession()
        },

        completeSession: { sessionId in
            SessionAPIClient.mockSession(id: sessionId, completed: true)
        },
        
        startRound: { sessionId in
            SessionAPIClient.mockRound(sessionId: sessionId, order: 1)
        },
        
        completeRound: { roundId, focusScore in
            SessionAPIClient.mockRoundCompletion(roundId: roundId, focusScore: focusScore)
        },

        getSession: { sessionId in
            SessionAPIClient.mockSession(id: sessionId)
        }
    )
}

extension DependencyValues {
    var sessionAPIClient: SessionAPIClient {
        get { self[SessionAPIClientKey.self] }
        set { self[SessionAPIClientKey.self] = newValue }
    }

    private enum SessionAPIClientKey: DependencyKey {
        static let liveValue = SessionAPIClient.live
    }
}
