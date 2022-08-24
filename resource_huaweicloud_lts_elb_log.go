package lts

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLtsElb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsElbCreate,
		ReadContext:   resourceLtsElbRead,
		DeleteContext: resourceLtsElbDelete,
		UpdateContext: resourceLtsElbUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLtsElbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["elb"], "https//", "https://", -1) + "v3/" + 
	       config.HwClient.ProjectID + "/elb/logtanks"
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	LogTank := entity.CreateLogTankOption{
		LogGroupId:     d.Get("log_group_id").(string),
		LoadBalancerId: d.Get("loadbalancer_id").(string),
		LogTopicId:     d.Get("log_topic_id").(string),
	}
	LogTankRequest := entity.CreateLogtankRequestBody{
		Logtank: &LogTank,
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrl(url).WithHeader(header).WithBody(LogTankRequest)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error creating LogTank fields %s : %s", LogTankRequest, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s, %s", string(body), err)
	}
	if response.StatusCode == 201 {
		rlt := &entity.CreateLogtankResponse{}
		err = json.Unmarshal(body, rlt)
		d.SetId(rlt.Logtank.ID)
		return resourceLtsElbRead(ctx, d, meta)
	}
	return diag.Errorf("error creating LogTank fields %s: %s", LogTankRequest, string(body))
}

func resourceLtsElbRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["elb"], "https//", "https://", -1) + "v3/" + config.HwClient.ProjectID +
	    "/elb/logtanks/" + d.Id()
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodGet).WithUrl(url).WithHeader(header)
	response, err := client.Do()
	body, diag := client.CheckDeletedDiag(d, err, response, "error Elb LogTank read instance")
	if body == nil {
		return diag
	}
	rlt := &entity.CreateLogtankResponse{}
	err = json.Unmarshal(body, rlt)
	mErr := multierror.Append(nil,
		d.Set("loadbalancer_id", rlt.Logtank.LoadBalancerID),
		d.Set("log_group_id", rlt.Logtank.LogGroupID),
		d.Set("log_group_id", rlt.Logtank.LogTopicID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting Elb LogTank fields: %w", err)
	}
	return nil
}

func resourceLtsElbDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["elb"], "https//", "https://", -1) + "v3/" + config.HwClient.ProjectID +
	    "/elb/logtanks/" + d.Id()
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"

	client.WithMethod(httpclient_go.MethodDelete).WithUrl(url).WithHeader(header)
	resp, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete LogTank %s: %s", d.Id(), err)
	}
	if resp.StatusCode == 204 {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error delete LogTank %s: %s", d.Id(), err)
	}
	return diag.Errorf("error delete LogTank %s:  %s", d.Id(), string(body))
}

func resourceLtsElbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["elb"], "https//", "https://", -1) + "v3/" + config.HwClient.ProjectID +
	    "/elb/logtanks/" + d.Id()
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	LogTankRequest := entity.CreateLogTankOption{
		LogGroupId: d.Get("log_group_id").(string),
		LogTopicId: d.Get("log_topic_id").(string),
	}
	client.WithMethod(httpclient_go.MethodPut).WithUrl(url).WithHeader(header).WithBody(LogTankRequest)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update LogTank fields %s: %s", LogTankRequest, err)
	}

	if response.StatusCode == 200 {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error update LogTank %s: %s", d.Id(), err)
	}
	return diag.Errorf("error update LogTank %s: %s", d.Id(), string(body))
}
