TAG_NAME?=dev-local

build:
	gcloud builds submit . \
		--tag gcr.io/$(PERSONAL_GCP_PROJECT)/movies-app:$(TAG_NAME) \
		--gcs-log-dir gs://trenner-movies/build-logs

deploy:
	gcloud run deploy movies-app \
		--image gcr.io/$(PERSONAL_GCP_PROJECT)/movies-app:$(TAG_NAME) \
		--platform managed \
		--max-instances 1 \
		--service-account movies-app@$(PERSONAL_GCP_PROJECT).iam.gserviceaccount.com \
		--set-env-vars "GOOGLE_CLOUD_PROJECT=$(PERSONAL_GCP_PROJECT)" \
		--region us-central1