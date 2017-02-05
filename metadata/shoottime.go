package metadata

type ShootTimeMetaDataHandler struct{}

func (handler ShootTimeMetaDataHandler) GetName() string {
	return "ShootTimeMetaDataHandler"
}

func (handler ShootTimeMetaDataHandler) UpdateNeeded(ctx *MetaDataHandlerContext) bool {
	return ctx.media.ShootTime.IsZero()
}

func (handler ShootTimeMetaDataHandler) Handle(ctx *MetaDataHandlerContext) error {
	dt, err := ctx.exifData.DateTime()
	if err != nil {
		return err
	}
	ctx.media.ShootTime = dt
	return nil
}
