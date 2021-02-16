SHELL := zsh
deploy-thor:
	helm install thor thor --set env.gcp.credentials="$(base64 $(pwd)/../gcloud-config.json)"