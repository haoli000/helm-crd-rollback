package io.javaoperatorsdk.operator.sample.v1alpha2;

@com.fasterxml.jackson.annotation.JsonInclude(com.fasterxml.jackson.annotation.JsonInclude.Include.NON_NULL)
@com.fasterxml.jackson.annotation.JsonPropertyOrder({"addedField", "reconciledBy"})
@com.fasterxml.jackson.databind.annotation.JsonDeserialize(using = com.fasterxml.jackson.databind.JsonDeserializer.None.class)
@javax.annotation.processing.Generated("io.fabric8.java.generator.CRGeneratorRunner")
public class LeaderElectionStatus implements io.fabric8.kubernetes.api.model.KubernetesResource {

    @com.fasterxml.jackson.annotation.JsonProperty("reconciledBy")
    @com.fasterxml.jackson.annotation.JsonSetter(nulls = com.fasterxml.jackson.annotation.Nulls.SKIP)
    private java.util.List<String> reconciledBy;
    @com.fasterxml.jackson.annotation.JsonProperty("addedField")
    private String addedField;

    public java.util.List<String> getReconciledBy() {
        return reconciledBy;
    }

    public void setReconciledBy(java.util.List<String> reconciledBy) {
        this.reconciledBy = reconciledBy;
    }

    public String getAddedField() {
        return addedField;
    }

    public void setAddedField(String addedField) {
        this.addedField = addedField;
    }
}

