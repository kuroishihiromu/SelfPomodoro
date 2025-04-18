//
//  ProfileScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import SwiftUI

struct ProfileScreenView: View {
    @State var username = ""
    @State var fullName = ""
    @State var website = ""

    @State var isLoading = false

    var body: some View {
        NavigationStack {
            Form {
                Section {
                    TextField("Username", text: $username)
                        .textContentType(.username)
                        .textInputAutocapitalization(.never)
                    TextField("Full name", text: $fullName)
                        .textContentType(.name)
                    TextField("Website", text: $website)
                        .textContentType(.URL)
                        .textInputAutocapitalization(.never)
                }

                Section {
                    Button("Update profile") {
                        updateProfileButtonTapped()
                    }
                    .bold()

                    if isLoading {
                        ProgressView()
                    }
                }
            }
            .navigationTitle("Profile")
            .toolbar(content: {
                ToolbarItem(placement: .topBarLeading) {
                    Button("Sign out", role: .destructive) {
                        Task {
                            try? await supabase.auth.signOut()
                        }
                    }
                }
            })
        }
        .task {
            await getInitialProfile()
        }
    }

    func getInitialProfile() async {
        do {
            let currentUser = try await supabase.auth.session.user

            let profile: Profile =
                try await supabase
                    .from("profiles")
                    .select()
                    .eq("id", value: currentUser.id)
                    .single()
                    .execute()
                    .value

            username = profile.username ?? ""
            fullName = profile.fullName ?? ""
            website = profile.website ?? ""

        } catch {
            debugPrint(error)
        }
    }

    func updateProfileButtonTapped() {
        Task {
            isLoading = true
            defer { isLoading = false }
            do {
                let currentUser = try await supabase.auth.session.user

                try await supabase
                    .from("profiles")
                    .update(
                        UpdateProfileParams(
                            username: username,
                            fullName: fullName,
                            website: website
                        )
                    )
                    .eq("id", value: currentUser.id)
                    .execute()
            } catch {
                debugPrint(error)
            }
        }
    }
}

#Preview {
    ProfileScreenView()
}
