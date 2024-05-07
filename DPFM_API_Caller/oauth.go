package dpfm_api_caller

import (
	dpfm_api_input_reader "data-platform-api-google-account-user-info-requests-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-google-account-user-info-requests-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-google-account-user-info-requests-rmq-kube/config"
	"encoding/json"
	"fmt"
	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
)

func (c *DPFMAPICaller) GoogleAccountUserInfo(
	input *dpfm_api_input_reader.SDC,
	errs *[]error,
	log *logger.Logger,
	conf *config.Conf,
) *[]dpfm_api_output_formatter.GoogleAccountUserInfoResponse {
	var googleAccountUserInfo []dpfm_api_output_formatter.GoogleAccountUserInfoResponse

	accessToken := input.GoogleAccountUserInfo.AccessToken

	userInfoURL := conf.OAuth.UserInfoURL

	req, err := http.NewRequest("GET", userInfoURL, nil)

	if err != nil {
		*errs = append(*errs, xerrors.Errorf("NewRequest error: %d", err))
		return nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		*errs = append(*errs, xerrors.Errorf("User info request error: %d", err))
		return nil
	}
	defer resp.Body.Close()

	userInfoBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		*errs = append(*errs, xerrors.Errorf("User info request response read error: %d", err))
		return nil
	}

	var response map[string]interface{}
	err = json.Unmarshal(userInfoBody, &response)
	if err != nil {
		*errs = append(*errs, xerrors.Errorf("Response response error: %d", err))
		return nil
	}

	errorObj, ok := response["error"].(map[string]interface{})
	if ok {
		code, ok := errorObj["code"].(float64)
		if ok {
			errMsg, _ := errorObj["message"].(string)
			*errs = append(*errs, xerrors.Errorf("Status code error: %v %v", code, errMsg))
			return nil
		}
	}

	var googleAccountUserInfoResponseBody dpfm_api_output_formatter.GoogleAccountUserInfoResponseBody
	err = json.Unmarshal(userInfoBody, &googleAccountUserInfoResponseBody)
	if err != nil {
		*errs = append(*errs, xerrors.Errorf("User info request response unmarshal error: %d", err))
		return nil
	}

	userInfo := dpfm_api_output_formatter.ConvertToGoogleAccountUserInfoRequestsFromResponse(googleAccountUserInfoResponseBody)

	googleAccountUserInfo = append(
		googleAccountUserInfo,
		userInfo,
	)

	return &googleAccountUserInfo
}
