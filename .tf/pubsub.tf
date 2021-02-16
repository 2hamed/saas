# resource "google_pubsub_topic" "primary" {
#   name = "${var.project_id}_capture_requests"
# }

# resource "google_pubsub_subscription" "primary" {
#   name  = "${var.project_id}_capture_requests"
#   topic = google_pubsub_topic.primary.name

#   # 20 minutes
#   message_retention_duration = "1200s"
#   retain_acked_messages      = true

#   ack_deadline_seconds = 20

#   expiration_policy {
#     ttl = "300000.5s"
#   }
#   retry_policy {
#     minimum_backoff = "10s"
#   }

#   enable_message_ordering = false
# }
