//
//  Supabase.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/26.
//

import Foundation
import Supabase

let supabase = SupabaseClient(
    supabaseURL: URL(string: "https://rnujnyggqrghhafeuqop.supabase.co")!,
    supabaseKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJudWpueWdncXJnaGhhZmV1cW9wIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDA1OTMyNzUsImV4cCI6MjA1NjE2OTI3NX0.d9Jg4jIdHs_HZqK5-v89rWRKE8XHF4hk10kbuJZWCTg"
)
