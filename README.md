https://api-console.zoho.eu/
ZohoCRM.modules.ALL,ZohoSearch.securesearch.READ,ZohoCRM.coql.READ,ZohoCRM.send_mail.all.CREATE,ZohoCRM.settings.emails.READ,ZohoCRM.settings.variables.ALL

curl -X POST -F grant_type=authorization_code -F client_id=1000.XXXXXXXXXX -F client_secret=XXXXX -F code=1000.XXXXXXXXXXXXXXXXXX https://accounts.zoho.eu/oauth/v2/token

curl -X POST "https://accounts.zoho.eu/oauth/v2/token?refresh_token=1000.XXXXXXXXXXXXXX&client_id=1000.XXXXXXXXXXXXX&client_secret=XXXXXXXXXX&grant_type=refresh_token"
