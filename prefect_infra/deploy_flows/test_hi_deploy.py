from prefect import flow
from prefect.deployments import Deployment

@flow(log_prints=True)
def test_hi():
    print("Hi from Prefect! 🤗")

def deploy():
    deployment = Deployment.build_from_flow(
        flow=test_hi,
        name="prefect-example-deployment"
    )
    deployment.apply()

if __name__ == "__main__":
    deploy()
