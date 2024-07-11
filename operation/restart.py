
import subprocess
import time

dc = '~/.kube/config.ma2'
app = 'haproxy-'
batch_size = 4
kubectl_cmd_prefix = "kubectl"



ingress_pods = []
# load pods
for line in subprocess.check_output(['{} --kubeconfig {} -n ingress get pods | grep {} | sort'.format(kubectl_cmd_prefix, dc,app).encode('utf-8')], shell=True, encoding='utf-8').split('\n'):
    if line != '':
        ingress_pods.append(line.split(" ")[0])

# batch_id = 0
batch = []
batches = []
for i in range(len(ingress_pods)):
    batch.append(ingress_pods[i])
    if len(batch) == batch_size:
        batches.append(batch)
        batch = []

if len(batch) > 0:
    batches.append(batch)

# let's restart stuff!
def is_ready():
    n_creating = 0
    for line in subprocess.check_output(['{} --kubeconfig {} -n ingress get pods | grep {}'.format(kubectl_cmd_prefix, dc,app).encode('utf-8')], shell=True, encoding='utf-8').split('\n'):
        if ("ContainerCreating" in line or "Terminating" in line or "1/2" in line):
            n_creating = n_creating+1
    print("{} containers still creating".format(n_creating))
    if(n_creating == 0):
        return True
    return False


for i in range(len(batches)):
    batch = batches[i]
    print("Batch {}:\n{}".format(i+1, batch))
    print("{} --kubeconfig {} -ningress delete pod {}".format(kubectl_cmd_prefix, dc, ' '.join(batch)))

for i in range(len(batches)):
    batch = batches[i]
    while not is_ready():
        time.sleep(10)
    print("k8s is ready; sleeping before starting next batch")
    time.sleep(10)
    cmd = "{} --kubeconfig {} -ningress delete pod {}".format(kubectl_cmd_prefix,dc, ' '.join(batch))
    print("Batch {} executing: {}".format(i+1, cmd))
    results = subprocess.check_output([cmd.encode('utf-8')], shell=True, encoding='utf-8').split('\n')
    print("Results: {}".format(results))
    
