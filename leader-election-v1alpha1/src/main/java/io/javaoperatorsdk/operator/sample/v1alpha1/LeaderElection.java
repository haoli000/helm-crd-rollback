package io.javaoperatorsdk.operator.sample.v1alpha1;

@io.fabric8.kubernetes.model.annotation.Version(value = "v1alpha1" , storage = true , served = true)
@io.fabric8.kubernetes.model.annotation.Group("sample.operator.javaoperatorsdk.io")
@io.fabric8.kubernetes.model.annotation.Singular("leaderelection")
@io.fabric8.kubernetes.model.annotation.Plural("leaderelections")
@javax.annotation.processing.Generated("io.fabric8.java.generator.CRGeneratorRunner")
public class LeaderElection extends io.fabric8.kubernetes.client.CustomResource<java.lang.Void, io.javaoperatorsdk.operator.sample.v1alpha1.LeaderElectionStatus> implements io.fabric8.kubernetes.api.model.Namespaced {
}

