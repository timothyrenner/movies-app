build:
	gcloud builds submit --tag gcr.io/$(PERSONAL_GCP_PROJECT)/movies-app

deploy:
	gcloud run deploy movies-app \
		--image gcr.io/$(PERSONAL_GCP_PROJECT)/movies-app \
		--platform managed \
		--max-instances 1 \
		--service-account movies-app@$(PERSONAL_GCP_PROJECT).iam.gserviceaccount.com \
		--set-env-vars "GOOGLE_CLOUD_PROJECT=$(PERSONAL_GCP_PROJECT)" \
		--region us-central1