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

func ResourceLtsStruct() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsStructTemplateCreate,
		ReadContext:   resourceLtsStructTemplateRead,
		DeleteContext: resourceLtsStructTemplateDelete,
		UpdateContext: resourceLtsStructTemplateUpdate,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parse_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regex_rules": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"layers": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tokenizer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_format": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"demo_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_analysis": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"field_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_defined_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"tag_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_analysis": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"demo_log": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// func buildDemoFieldsInfo(rawDimensions []interface{}) []entity.DemoFieldsInfo {
// 	if len(rawDimensions) == 0 {
// 		return nil
// 	}
// 	dimensions := make([]entity.DemoFieldsInfo, len(rawDimensions))
// 	for i, rawdimension := range rawDimensions {
// 		dimension := rawdimension.(map[string]interface{})
// 		dimensions[i] = entity.DemoFieldsInfo{
// 			IsAnalysis:      dimension["is_analysis"].(bool),
// 			Content:         dimension["content"].(string),
// 			FieldName:       dimension["field_name"].(string),
// 			Type:            dimension["type"].(string),
// 			UserDefinedName: dimension["user_defined_name"].(string),
// 			Index:           dimension["index"].(int),
// 		}
// 	}
// 	return dimensions
// }

// func buildTagFieldsInfo(rawDimensions []interface{}) []entity.TagFieldsInfo {
// 	if len(rawDimensions) == 0 {
// 		return nil
// 	}
// 	dimensions := make([]entity.TagFieldsInfo, len(rawDimensions))
// 	for i, rawdimension := range rawDimensions {
// 		dimension := rawdimension.(map[string]interface{})
// 		dimensions[i] = entity.TagFieldsInfo{
// 			FieldName:  dimension["field_name"].(string),
// 			Type:       dimension["type"].(string),
// 			Content:    dimension["content"].(*string),
// 			IsAnalysis: dimension["is_analysis"].(*bool),
// 		}
// 	}
// 	return dimensions
// }
func resourceLtsStructTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	var url string
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	opts := entity.StructTemplateRequest{
		LogGroupId:  d.Get("log_group_id").(string),
		LogStreamId: d.Get("log_stream_id").(string),
	}
	if d.Get("template_type").(string) == "custom" {
		url = strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v2/" + 
	       config.HwClient.ProjectID + "/lts/struct/template"
		opts.ToDemoFieldsInfo()
		opts.ParseType = "split"
		opts.Tokenizer = " "
		opts.Content = "127.0.0.1 10.142.203.101 8080 [18/Aug/2021:15:14:33 +0800] GET /apm HTTP/1.1 404 86 6"
	} else {
		url = strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v3/" + 
			config.HwClient.ProjectID + "/lts/struct/template"
			opts.TemplateId = d.Get("template_id").(string)
			opts.TemplateType = d.Get("template_type").(string)
			opts.TemplateName = d.Get("template_name").(string)
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrl(url).WithHeader(header).WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error request creating StructTemplate fields %s: %s", opts, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s , %s", string(body), err)
	}
	if response.StatusCode == 201 || response.StatusCode == 200 {
		return resourceLtsStructTemplateRead(ctx, d, meta)
	}
	return diag.Errorf("error creating StructTemplate fields %s: %s", opts, string(body))
}

func resourceLtsStructTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v2/" +
	    config.HwClient.ProjectID + "/lts/struct/template?logGroupId=" +
		d.Get("log_group_id").(string) + "&logStreamId=" + d.Get("log_stream_id").(string)
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	client.WithMethod(httpclient_go.MethodGet).WithUrl(url).WithHeader(header)
	resp, err := client.Do()
	body, diag := client.CheckDeletedDiag(d, err, resp, "error StructTemplate read instance")
	if body == nil {
		return diag
	}
	body = body[1 : len(body)-1]
	body2 := strings.Replace(string(body), `\\\`, "**", -1)
	body3 := strings.Replace(body2, `\`, "", -1)
	body4 := strings.Replace(body3, "**", `\`, -1)
	rlt := &entity.ShowStructTemplateResponse{}
	err = json.Unmarshal([]byte(body4), rlt)
	d.SetId(rlt.Id)
	mErr := multierror.Append(nil,
		d.Set("demo_log", rlt.DemoLog),
		d.Set("log_group_id", rlt.LogGroupId),
		d.Set("log_stream_id", rlt.LogStreamId),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting LtsStructTemplate fields: %w", err)
	}
	return nil
}

func resourceLtsStructTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v2/" + 
	    config.HwClient.ProjectID + "/lts/struct/template"
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	structTemplateDeleteRequest := entity.DeleteStructTemplateReqBody{
		Id: d.Id(),
	}
	client.WithMethod(httpclient_go.MethodDelete).WithUrl(url).WithHeader(header).WithBody(structTemplateDeleteRequest)
	resp, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete StructTemplate %s: %s", d.Id(), err)
	}
	if resp.StatusCode == 200 {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error delete StructTemplate %s: %s", d.Id(), string(body))
	}
	return nil
}

func resourceLtsStructTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(config)
	if diaErr != nil {
		return diaErr
	}
	url := strings.Replace(config.Endpoints["lts"], "https//", "https://", -1) + "v2/" + 
	    config.HwClient.ProjectID + "/lts/struct/template"
	header := make(map[string]string)
	header["content-type"] = "application/json;charset=UTF8"
	structTemplateRequest := entity.StructTemplateRequest{
		LogGroupId:   d.Get("log_group_id").(string),
		LogStreamId:  d.Get("log_stream_id").(string),
		TemplateId:   d.Get("template_id").(string),
		TemplateType: d.Get("template_type").(string),
		TemplateName: d.Get("template_name").(string),
	}
	client.WithMethod(httpclient_go.MethodPut).WithUrl(url).WithHeader(header).WithBody(structTemplateRequest)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update StructTemplate fields %s: %s", structTemplateRequest, err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s , %s", string(body), err)
	}
	if response.StatusCode == 201 {
		rlt := &entity.DeleteStructTemplateReqBody{}
		err = json.Unmarshal(body, rlt)
		d.SetId(rlt.Id)
		return nil
	}
	return diag.Errorf("error update StructTemplate fields %s: %s", structTemplateRequest, err)
}
