docker build -t practice-operator:0.X .
docker tag  practice-operator:0.X jackfantasy/practice-operator:0.X
docker push jackfantasy/practice-operator:0.X

and then change the image inside the practice-operator/artifacts/practice-operator.yaml
