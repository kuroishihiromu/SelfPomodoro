// //
// //  RootContentView.swift
// //  SelfPomodoro
// //
// //  Created by tsunakit99 on 2025/11/05.
// //

// import SwiftUI
// import Amplify
// import AWSCognitoAuthPlugin
// import AWSCognitoIdentityProvider
// import AWSPluginsCore
// import AWSAPIPlugin
// import ComposableArchitecture

// struct _OldContentView: View {
//     @State private var isConfigured = false
//     @State private var isSignedIn: Bool? = nil
//     @State private var tokens: AuthTokens? = nil

//     var body: some View {
//         ZStack {
//             if !isConfigured {
//                 ProgressView("Initializing...")
//             } else if isSignedIn == true, let tokens {
//                 AuthScreenView(
//                     store: Store(
//                         initialState: AuthFeature.State(
//                             tokens: tokens,
//                             isLoggedIn: true
//                         ),
//                         reducer: { AuthFeature() }
//                     )
//                 )
//             } else {
//                 AuthScreenView(
//                     store: Store(
//                         initialState: AuthFeature.State(),
//                         reducer: { AuthFeature() }
//                     )
//                 )
//             }
//         }
//         .task {
//             await initializeAppIfNeeded()
//         }
//     }

//     @MainActor
//     private func initializeAppIfNeeded() async {
//         guard !isConfigured else { return }

//         do {
//             try Amplify.add(plugin: AWSCognitoAuthPlugin())
//             try Amplify.add(plugin: AWSAPIPlugin())
//             try Amplify.configure()
//             isConfigured = true

//             let session = try await Amplify.Auth.fetchAuthSession()
//             if session.isSignedIn,
//                let provider = session as? AuthCognitoTokensProvider {
//                 let result = provider.getCognitoTokens()
//                 let token = try result.get()
//                 let authTokens = AuthTokens(
//                     idToken: token.idToken,
//                     accessToken: token.accessToken,
//                     refreshToken: token.refreshToken
//                 )
//                 tokens = authTokens
//                 isSignedIn = true
//             } else {
//                 isSignedIn = false
//             }
//         } catch {
//             if let authError = error as? AuthError,
//                case .service(_, _, let underlyingError) = authError,
//                let cognitoError = underlyingError as? AWSCognitoIdentityProvider.NotAuthorizedException,
//                cognitoError.properties.message?.contains("Refresh Token has expired") == true {
//                 // Refresh token expired: treat as signed out
//             }
//             isSignedIn = false
//         }
//     }
// }
