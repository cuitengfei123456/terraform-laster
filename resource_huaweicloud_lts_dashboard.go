package lts

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLtsDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsDashBoardCreate,
		ReadContext:   resourceLtsDashBoardRead,
		DeleteContext: resourceLtsDashBoardDelete,
		UpdateContext: resourceDashBoardUpdate,
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
			"is_delete_charts": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_stream_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_title": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"template_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"last_update_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},	
		},
	}
}

func resourceLtsDashBoardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v2/" + 
	       config.HwClient.ProjectID + "/lts/template-dashboard"
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	dashBoardRequest := entity.DashBoardRequest{
		LogGroupId:    d.Get("log_group_id").(string),
		LogGroupName:  d.Get("log_group_name").(string),
		LogStreamId:   d.Get("log_stream_id").(string),
		LogStreamName: d.Get("log_stream_name").(string),
		TemplateTitle: utils.ExpandToStringList(d.Get("template_title").([]interface{})),
		TemplateType:  utils.ExpandToStringList(d.Get("template_type").([]interface{})),
		GroupName:     d.Get("group_name").(string),
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrl(url).WithHeader(header).WithBody(dashBoardRequest)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error creating LtsDashBoard fields %s: %s", dashBoardRequest, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s, %s", string(body), err)
	}
	if response.StatusCode == 201 {
		rlt := make([]entity.DashBoard, 0)
		err = json.Unmarshal(body, &rlt)
		if err != nil {
			return diag.Errorf("error convert data %s , %s", string(body), err)
		}
		if len(rlt) == 0 {
			return diag.Errorf("error resource has been created log stream name %s", d.Get("log_stream_name").(string))
		}
		d.SetId(rlt[0].Id)
		return  resourceLtsDashBoardRead(ctx, d, meta)
	}
	return diag.Errorf("error creating LtsDashBoard Response %s: %s", dashBoardRequest, string(body))
}

func resourceLtsDashBoardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v2/" + 
	    config.HwClient.ProjectID + "/dashboards?id=" + d.Id()
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodGet).WithUrl(url).WithHeader(header)
	response, err := client.Do()
	body, diags := client.CheckDeletedDiag(d, err, response, "error LtsDashBoard")
	if body == nil {
		return diags
	}
	rlt := entity.ReadDashBoardResp{}
	err = json.Unmarshal(body, &rlt)
	d.Set("region", config.GetRegion(d))
	if err != nil || len(rlt.Results) == 0 {
		return diag.Errorf("error read lts dash board %s", d.Id())
	}
	mErr := multierror.Append(nil,
		d.Set("title", rlt.Results[0].Title),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting Lts dashboard fields: %s", err)
	}
	return nil
}

func resourceLtsDashBoardDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v2/" + 
	    config.HwClient.ProjectID + "/dashboard?is_delete_charts=" +d.Get("is_delete_charts").(string) + "&id=" + d.Id()
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodDelete).WithUrl(url).WithHeader(header)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete LtsDashBoard %s: %s", d.Id(), err)
	}
	if response.StatusCode == 200 {
		return nil
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error delete LtsDashBoard %s: %s", d.Id(), err)
	}
	return diag.Errorf("error delete LtsDashBoard %s:  %s", d.Id(), string(body))
}


func resourceDashBoardUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v2/" + 
	    config.HwClient.ProjectID + "/dashboard?id=" + d.Id()
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	dashBoardRequest := entity.DashBoardRequest{
		LogGroupId:    d.Get("log_group_id").(string),
		LogGroupName:  d.Get("log_group_name").(string),
		LogStreamId:   d.Get("log_stream_id").(string),
		LogStreamName: d.Get("log_stream_name").(string),
		TemplateTitle: utils.ExpandToStringList(d.Get("template_title").([]interface{})),
		TemplateType:  utils.ExpandToStringList(d.Get("template_type").([]interface{})),
		GroupName:     d.Get("group_name").(string),
	}
	client.WithMethod(httpclient_go.MethodPut).WithUrl(url).WithHeader(header).WithBody(dashBoardRequest)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update LtsDashBoard fields %s: %s", dashBoardRequest, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s: %s", string(body), err)
	}
	if response.StatusCode == 200 {
		rlt := make([]entity.DashBoard, 0)
		err = json.Unmarshal(body, &rlt)
		if err != nil {
			return diag.Errorf("error convert data %s: %s", string(body), err)
		}
		d.SetId(rlt[0].Id)
		return nil
	}
	return diag.Errorf("error update LtsDashBoard fields %s: %s", dashBoardRequest, string(body)) 
}
