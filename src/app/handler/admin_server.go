package handler

import (
	pb "github.com/kujilabo/cocotola-translator-api/src/proto"
)

type adminServer struct {
	pb.UnimplementedTranslatorAdminServer
}
