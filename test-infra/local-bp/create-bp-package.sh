REPO='oai-ran-bp'
echo "This script is only for testing"
mkdir temp
cd temp

mkdir master
git clone https://github.com/jain-ashish-sam/catalog.git
mv catalog/workloads/oai/* master/
rm -rf catalog
cp ../deployment.yaml master/oai-ran-operator/operator/

# for pkg in 'pkg-example-cucp-bp' 'pkg-example-cuup-bp' 'pkg-example-du-bp';
for pkg in 'oai-ran-operator' 'pkg-example-cucp-bp' 'pkg-example-cuup-bp' 'pkg-example-du-bp';
do
    CREATED_PKG=$(kpt alpha rpkg init --repository=$REPO $pkg --workspace=v1 -ndefault| awk '{print $1;}')
    kpt alpha rpkg pull  $CREATED_PKG ./yourpkg -ndefault
    cp -r master/$pkg/*  yourpkg/
    echo "Editing Package $pkg is Done| Now Push-propose-Approve"
    kpt alpha rpkg push $CREATED_PKG yourpkg -ndefault
    kpt alpha rpkg propose $CREATED_PKG -ndefault
    kpt alpha rpkg approve $CREATED_PKG -ndefault
    rm -rf yourpkg
done

cd ..
rm -rf temp
echo "Done| Check the content of your blueprint-repo once"


