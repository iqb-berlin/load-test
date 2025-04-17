# Preparation

- Configure test parameters in config.json
    (- Use the url-list-gen.py script in the root of the project to generate the list)

*When testing the Testcenter version with file service, you need to change the path for resource files.
Swap the commented line in the script.*

# Running Locust
- First install all dependencies with `pip install -r requirements.txt`
- Run `locust` and open the UI on Port 8089

*When testing a lot of users you may want to run several Locust-workers. Those can run on the same machine.
Start and put them in the background with `locust --worker &`. Do this once per available core.
