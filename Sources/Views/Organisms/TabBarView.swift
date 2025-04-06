//
//  TabBarView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI
struct TabBarView: View {
    var body: some View {
        TabView {
            HomeScreenView()
                .tabItem {
                    Image(systemName: "house")
                }

            TaskManagementScreenView()
                .tabItem {
                    Image(systemName: "checkmark.circle")
                }

            StatisticsScreenView()
                .tabItem {
                    Image(systemName: "chart.bar")
                }

            ProfileScreenView()
                .tabItem {
                    Image(systemName: "person.crop.circle")

                }
        }
    }
}

#Preview {
    TabBarView()
}
