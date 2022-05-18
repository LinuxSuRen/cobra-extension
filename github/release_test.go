package github_test

import (
	jClient "github.com/linuxsuren/cobra-extension/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	ghClient := jClient.ReleaseClient{}

	assert.Nil(t, ghClient.Client)
	ghClient.Init()
	assert.NotNil(t, ghClient.Client)
}

func TestGetLatestReleaseAsset(t *testing.T) {
	client, teardown := jClient.PrepareForGetLatestReleaseAsset() //setup()
	defer teardown()

	ghClient := jClient.ReleaseClient{
		Client: client,
	}
	asset, err := ghClient.GetLatestReleaseAsset("o", "r")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

func TestGetReleaseAssetByTagName(t *testing.T) {
	client, teardown := jClient.PrepareForGetReleaseAssetByTagName() //setup()
	defer teardown()

	ghClient := jClient.ReleaseClient{
		Client: client,
	}
	asset, err := ghClient.GetReleaseAssetByTagName("jenkins-zh", "jenkins-cli", "tagName")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}
