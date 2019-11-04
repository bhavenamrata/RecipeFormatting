# Building TensorFlow

The instructions provided below specify the steps to build Tensorflow version 2.0.0 on Linux on IBM Z for

* Ubuntu (16.04, 18.04, 19.04)

### _**General Notes:**_  
* _When following the steps below please use a standard permission user unless otherwise specified._
* _A directory `/<source_root>/` will be referred to in these instructions, this is a temporary writable directory anywhere you'd like to place it._


## Step 1: Building and Installing TensorFlow v2.0.0

#### 1.1) Build using script

If you want to build Tensorflow using manual steps, go to STEP 1.2.

Use the following commands to build Tensorflow using the build [script](https://github.com/linux-on-ibm-z/scripts/tree/master/Tensorflow). Please make sure you have wget installed.

```
wget -q https://raw.githubusercontent.com/linux-on-ibm-z/scripts/master/Tensorflow/2.0.0/build_tensorflow.sh

# Build Tensorflow
bash build_tensorflow.sh    [Provide -t option for executing build with tests]
```

If the build completes successfully, go to STEP 2. In case of error, check `logs` for more details or go to STEP 1.2 to follow manual build steps.
 
#### 1.2) Install the dependencies

 ```shell  
 export SOURCE_ROOT=/<source_root>/
 ```  
 

  * Ubuntu 16.04
 ```shell
  sudo apt-get update  
  sudo apt-get install -y pkg-config zip g++ zlib1g-dev unzip git vim tar wget automake autoconf libtool make curl maven openjdk-11-jdk python3-pip python3-virtualenv python3-numpy swig python3-dev libcurl3-dev python3-mock python3-scipy bzip2 python3-sklearn libhdf5-dev patch git patch libssl-dev
  sudo pip3 install numpy==1.16.2 future wheel backports.weakref portpicker futures==2.2.0 enum34 keras_preprocessing keras_applications h5py tensorflow_estimator
  
 ```
	
  * (For Ubuntu 16.04) Download and install AdoptOpenJDK (OpenJDK 11 with HotSpot) from [here](https://adoptopenjdk.net/releases.html?variant=openjdk11&jvmVariant=hotspot#s390x_linux)  
 ```shell
  export JAVA_HOME=/<path to JDK>/
  export PATH=$JAVA_HOME/bin:$PATH
 ```
 
  * Ubuntu (18.04, 19.04)
 ```shell
  sudo apt-get update  
  sudo apt-get install -y pkg-config zip g++ zlib1g-dev unzip git vim tar wget automake autoconf libtool make curl maven openjdk-11-jdk python3-pip python3-virtualenv python3-numpy swig python3-dev libcurl3-dev python3-mock python3-scipy bzip2 python3-sklearn libhdf5-dev patch git patch libssl-dev
  sudo pip3 install numpy==1.16.2 future wheel backports.weakref portpicker futures enum34 keras_preprocessing keras_applications h5py tensorflow_estimator
  
 ```
 
 
  * Install grpcio   
 ```shell  
  export GRPC_PYTHON_BUILD_SYSTEM_OPENSSL=True
  sudo -E pip3 install grpcio  
 ```  
  
  * Install go   
 ```shell  
  cd $SOURCE_ROOT   
  wget https://dl.google.com/go/go1.13.3.linux-s390x.tar.gz  
  tar -C /usr/local -xzf go1.13.3.linux-s390x.tar.gz  
  export PATH=/usr/local/go/bin:$PATH  
  go version   
 ```    
  
#### 1.3)  Build Bazel

* Download Bazel  
  ```shell  
   cd $SOURCE_ROOT   
   mkdir bazel && cd bazel  
   wget https://github.com/bazelbuild/bazel/releases/download/0.26.1/bazel-0.26.1-dist.zip
   unzip bazel-0.26.1-dist.zip 
   chmod -R +w .
  ```  

* Build Bazel  

  ```shell  

  env EXTRA_BAZEL_ARGS="--host_javabase=@local_jdk//:jdk" bash ./compile.sh
  export PATH=$PATH:$SOURCE_ROOT/bazel/output/  
  ```  
  
  _**Note:** While building Bazel, if build fails with an error `java.lang.OutOfMemoryError: Java heap space`, apply below patch and rebuild Bazel._  
  
  * Create a patch file `$SOURCE_ROOT/bazel/scripts/bootstrap/patch_compile.diff` with the following contents:
  
     ```diff  
     @@ -127,7 +127,7 @@ function java_compilation() {
         # Useful if your system chooses too small of a max heap for javac.
         # We intentionally rely on shell word splitting to allow multiple
         # additional arguments to be passed to javac.

     -   run "${JAVAC}" -classpath "${classpath}" -sourcepath "${sourcepath}" \
     +   run "${JAVAC}" -J-Xms1g -J-Xmx1g -classpath "${classpath}" -sourcepath "${sourcepath}" \
       -d "${output}/classes" -source "$JAVA_VERSION" -target "$JAVA_VERSION" \
       -encoding UTF-8 ${BAZEL_JAVAC_OPTS} "@${paramfile}"

     ``` 

  * Apply the patch file  
     ```shell  
     cd $SOURCE_ROOT/bazel/scripts/bootstrap/
     patch compile.sh < patch_compile.diff 
     ```  

  * Rebuild 
     ```shell  
     cd $SOURCE_ROOT/bazel/
     env EXTRA_BAZEL_ARGS="--host_javabase=@local_jdk//:jdk" bash ./compile.sh
     export PATH=$PATH:$SOURCE_ROOT/bazel/output/
     ```   

#### 1.4)  Build TensorFlow

* Download source code
  ```shell
  cd $SOURCE_ROOT
  git clone https://github.com/linux-on-ibm-z/tensorflow
  cd tensorflow
  git checkout v2.0.0-s390x
  ```  

* Configure    

  ```shell
  ./configure  
  Extracting Bazel installation...
  You have bazel 0.26.1- (@non-git) installed.
  Please specify the location of python. [Default is /usr/bin/python]: */usr/bin/python3*

  Found possible Python library paths:
    /usr/lib/python3/dist-packages
    /usr/local/lib/python3.7/dist-packages
  Please input the desired Python library path to use.  Default is [/usr/lib/python3/dist-packages]
  
  Do you wish to build TensorFlow with XLA JIT support? [Y/n]: **N**
  No XLA JIT support will be enabled for TensorFlow.

  Do you wish to build TensorFlow with OpenCL SYCL support? [y/N]: N
  No OpenCL SYCL support will be enabled for TensorFlow.

  Do you wish to build TensorFlow with ROCm support? [y/N]: N
  No ROCm support will be enabled for TensorFlow.

  Do you wish to build TensorFlow with CUDA support? [y/N]: N
  No CUDA support will be enabled for TensorFlow.

  Do you wish to download a fresh release of clang? (Experimental) [y/N]: N
  Clang will not be downloaded.

  Do you wish to build TensorFlow with MPI support? [y/N]: N
  No MPI support will be enabled for TensorFlow.

  Please specify optimization flags to use during compilation when bazel option "--config=opt" is specified [Default is -march=native -Wno-sign-compare]:


  Would you like to interactively configure ./WORKSPACE for Android builds? [y/N]: N
  Not configuring the WORKSPACE for Android builds.

  Preconfigured Bazel build configs. You can use any of the below by adding "--config=<>" to your build command. See .bazelrc for more details.
        --config=mkl            # Build with MKL support.
        --config=monolithic     # Config for mostly static monolithic build.
        --config=gdr            # Build with GDR support.
        --config=verbs          # Build with libverbs support.
        --config=ngraph         # Build with Intel nGraph support.
        --config=numa           # Build with NUMA support.
        --config=dynamic_kernels        # (Experimental) Build kernels into separate shared objects.
        --config=v2             # Build TensorFlow 2.x instead of 1.x.
  Preconfigured Bazel build configs to DISABLE default on features:
        --config=noaws          # Disable AWS S3 filesystem support.
        --config=nogcp          # Disable GCP support.
        --config=nohdfs         # Disable HDFS support.
        --config=noignite       # Disable Apache Ignite support.
        --config=nokafka        # Disable Apache Kafka support.
        --config=nonccl         # Disable NVIDIA NCCL support.
  Configuration finished

  ```  

* Build Tensorflow

  ```shell
  bazel --host_jvm_args="-Xms1024m" --host_jvm_args="-Xmx2048m" build  --define=tensorflow_mkldnn_contraction_kernel=0 //tensorflow/tools/pip_package:build_pip_package
  ``` 

#### 1.5)  Build and install TensorFlow wheel

  ```shell  
  cd $SOURCE_ROOT/tensorflow
  bazel-bin/tensorflow/tools/pip_package/build_pip_package /tmp/tensorflow_wheel
  sudo pip3 install /tmp/tensorflow_wheel/tensorflow-2.0.0-cp*-linux_s390x.whl
  ```  

## Step 2: Verify TensorFlow (Optional)  
* Run TensorFlow from command Line   

  ```shell
   $ cd $SOURCE_ROOT
   $ python3
    >>> import tensorflow as tf
    >>> tf.add(1, 2).numpy()
    3
    >>> hello = tf.constant('Hello, TensorFlow!')
    >>> hello.numpy()
    'Hello, TensorFlow!'
    >>>  
  ```  

## Step 3: Execute Test Suite (Optional)  

* Run complete testsuite  

  ```shell
  bazel --host_jvm_args="-Xms1024m" --host_jvm_args="-Xmx2048m" test --define=tensorflow_mkldnn_contraction_kernel=0 --host_javabase="@local_jdk//:jdk" --test_tag_filters=-gpu,-benchmark-test,-v1only -k   --test_timeout 300,450,1200,3600 --build_tests_only --test_output=errors -- //tensorflow/... -//tensorflow/compiler/... -//tensorflow/lite/... -//tensorflow/core/platform/cloud/... -//tensorflow/java/... -//tensorflow/contrib/... 
  ```
  _**Note:** Skipping some test modules due to below issues:          
   *  `//tensorflow/lite` and `//tensorflow/core/platform/cloud` due to BoringSSL : `#error Unknown target CPU`, for more details refer [#14039](https://github.com/tensorflow/tensorflow/issues/14039) 
   *  `//tensorflow/java` due to error `Building Java resource jar failed ` for more details refer [#19770](https://github.com/tensorflow/tensorflow/issues/19770) 
   *  `//tensorflow/contrib` skipped as tf.contrib has been deprecated from v2.x onwards_      
  
  
* Run individual test 
  ```shell
  bazel --host_jvm_args="-Xms1024m" --host_jvm_args="-Xmx2048m" test --define=tensorflow_mkldnn_contraction_kernel=0 --host_javabase="@local_jdk//:jdk" //tensorflow/<module_name>:<testcase_name>
  ```  
    For example,   
    ```shell
    bazel --host_jvm_args="-Xms1024m" --host_jvm_args="-Xmx2048m" test --define=tensorflow_mkldnn_contraction_kernel=0 --host_javabase="@local_jdk//:jdk" //tensorflow/python/kernel_tests:topk_op_test
    ```  
 
  _**Note:**_       
  _1. Below tests are failing on s390x and those are either known or equivalent to Intel:_  
     `//tensorflow/core/grappler/costs:graph_properties_test`  
     `//tensorflow/python:file_io_test`  
     `//tensorflow/python:session_clusterspec_prop_test`  
     `//tensorflow/python/autograph/pyct:inspect_utils_test_par`  
     `//tensorflow/python/compiler/xla:xla_test`  
     `//tensorflow/python/debug:debugger_cli_common_test`  
     `//tensorflow/python/debug:dist_session_debug_grpc_test`  
     `//tensorflow/python/eager:backprop_test`  
     `//tensorflow/python/eager:def_function_xla_test_cpu`  
     `//tensorflow/python/kernel_tests:reader_ops_test`  
     `//tensorflow/python/ops/parallel_for:xla_control_flow_ops_test`  
     `//tensorflow/python/tpu:tpu_test`  
     `//tensorflow/python/training/tracking:util_xla_test_cpu`  
     `//tensorflow/python/tpu:datasets_test`  
	 `//tensorflow/python/kernel_tests/random:random_binomial_test`  
	 
  _2. Below tests are failing on s390x and investigation is in progress:_     
     `//tensorflow/python/kernel_tests:unicode_decode_op_test`  
     `//tensorflow/python/kernel_tests:unicode_transcode_op_test`  
     `//tensorflow/python:cluster_test`  
     `//tensorflow/python:cost_analyzer_test`  
  

## References:
   https://www.tensorflow.org/  
   https://github.com/tensorflow/tensorflow   
   http://bazel.io/  
   

