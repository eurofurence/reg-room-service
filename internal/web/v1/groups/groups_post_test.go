package groups

// import (
//	"context"
//	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
//	"github.com/eurofurence/reg-room-service/internal/web/common"
//	"net/http"
//	"net/http/httptest"
//)
//
//func (g *groupsSuite) TestCreateGroup() {
//	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://base.test.url/groups", nil)
//	g.NoError(err)
//
//	ctx := req.Context()
//	ctx = context.WithValue(ctx, common.CtxKeyRequestURL{}, req.URL)
//
//	recorder := httptest.NewRecorder()
//
//	_, err = g.ctrl.CreateGroup(ctx, &CreateGroupRequest{Group: modelsv1.GroupCreate{
//		Name:  "test-group",
//		Flags: []string{"flag1", "flag2"},
//	}}, recorder)
//
//	locationHeader := recorder.Header().Get("Location")
//
//	g.Equal("", locationHeader)
//}
//
