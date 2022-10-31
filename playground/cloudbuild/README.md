
# Setup of Cloud Build and GitHub repository to deploy Beam Playground

**This README explains how to:**

1. Enable Cloud Build in your GCP Project

2. Connect Github repository to Cloud Build

3. Set up Cloud Build service account (IAM)


### Before you begin

- Create GCP project
    - [Create a project docs](https://cloud.google.com/apis/docs/getting-started#creating_a_google_project)
- Enable billing for the project
    - [Enable billing docs](https://cloud.google.com/apis/docs/getting-started#enabling_billing)
- Enable Cloud Build API
    - [Enable API docs](https://cloud.google.com/apis/docs/getting-started#enabling_apis)
- Have your GitHub repository ready
- Have either a Dockerfile or a [Cloud Build config file](https://cloud.google.com/build/docs/build-config) in your GitHub source repository.
- If you are initially connecting your repository to Cloud Build, make sure you have admin-level permissions on your repository. To learn more about GitHub repository permissions, see [Repository permission levels for an organization](https://docs.github.com/en/github/setting-up-and-managing-organizations-and-teams/repository-permission-levels-for-an-organization#permission-levels-for-repositories-owned-by-an-organization)


## Enable Cloud Build and connect a GitHub Repository:

To connect your GitHub repository to your Cloud Build:

1. Open the Navigation Menu on the left side of your Google Cloud Platform Console.

2. Select Cloud Build, if it did not appear there, choose View All Products and search for Cloud Build there

3. Open the Triggers page in the Cloud Build console.

4. In the project selector in the top bar, select your Cloud project.

5. Click Connect repository (on the top).

6. Select GitHub (Cloud Build GitHub App), check the consent checkbox, and click Continue.

7. (Optional) If you have not signed into GitHub before, do so now.

        The Authorization page will appear where you are asked to authorize the Google Cloud Build App to connect to Google Cloud.
    
        Click Authorize Google Cloud Build by GoogleCloudBuild.

8. Click Install Google Cloud Build.

9. In the pop-up that appears, select your GitHub username or organization.

10. Select one of the following options based on your business need:

        a. All repositories - enable current and future GitHub repositories for access via the Cloud Build app

        b. Only select repositories - use the Select repositories drop-down to enable only specific repositories for access via the Cloud Build app.

        *You are able to enable additional repositories at a later time. If you select All repositories as your option, the Cloud Build app is authorized to access all your repositories. However, you need to connect each new repository through Cloud Build following the steps outlined in this section.*

11. Click Install to install the Cloud Build app.

        The pop-up closes and you are directed to a project selector page within Cloud Build. On this page, you can select an existing Cloud project or create a new project.

        If you do not see an existing project listed on this page, click Select Project to see a list of all existing projects.

12. After you have selected a project or created a new one, you will see the Connect repository panel.

13. In the Select repository section, select the following fields:

* GitHub account: The GitHub account used to install the Cloud Build GitHub App. This field may be pre-selected for you.

* Repository: The repositories you want to connect to Cloud Build.

         If you don't see one or more of your target repositories, click Edit repositories on GitHub and repeat the steps above to enable additional repositories in the Cloud Build GitHub App.

14. Once you have selected your GitHub account and repositories, read the consent disclaimer and select the checkbox next to it to indicate that you accept the presented terms.

15. Click Connect.

16. (Optional) In the Create a trigger section, select the repositories you want to create a trigger for in the Create a sample trigger for these repositories field. Once you have selected your repositories, click Create a trigger.

17. Click Done.

You have now connected one or more GitHub repositories to your Cloud project. You are directed to the Triggers page in Google Cloud console.

## Set up Cloud Build service account (IAM)

> Cloud Build executes builds on your behalf using a unique service account.
> The Cloud Build service account is automatically established and given the Cloud Build Service Account role for the project when the Cloud Build API is enabled on a Google Cloud project.
> This role grants the service account authority to carry out a number of actions, but you can offer the service account extra permissions to carry out more tasks.
> How to grant and withdraw rights for the Cloud Build service account are described on this page.

### Before you begin

- Understand [Cloud Build roles and permissions](https://cloud.google.com/build/docs/iam-roles-permissions).
- Read [Cloud Build service account](https://cloud.google.com/build/docs/cloud-build-service-account).


**Granting a role to the Cloud Build service account using the IAM page**


1. Open the Navigation Menu on the left side of your Google Cloud Platform Console.

2. Select IAM & Admin option.

3. Select your Cloud project on the top.

4. In the permissions table, locate the row with the email address ending with @cloudbuild.gserviceaccount.com. This is your Cloud Build service account.

5. Click on the pencil icon.

6. Select the next roles to grant to the Cloud Build service account:

        App Engine Admin
        App Engine Creator
        Artifact Registry Administrator
        Cloud Build Service Account
        Cloud Datastore Index Admin
        Cloud Memorystore Redis Admin
        Compute Admin
        Create Service Accounts
        Kubernetes Engine Admin
        Quota Administrator
        Role Administrator
        Security Admin
        Service Account User
        Storage Admin

7. Click Save.