# resource "google_spanner_instance" "primary" {
#   config       = "regional-us-central1"
#   display_name = var.project_id
#   num_nodes    = 1
# }

# resource "google_spanner_database" "database" {
#   instance = google_spanner_instance.primary.name
#   name     = var.project_id
#   ddl = [
#     "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
#     "CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
#   ]
#   deletion_protection = false
# }
