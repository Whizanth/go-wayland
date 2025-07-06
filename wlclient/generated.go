package wlclient

import "git.whizanth.com/go/wayland"

type Object struct {
	client *Client
	id     uint32
}

func (object Object) Id() uint32 {
	return object.id
}

func New() (*Client, error) {
	client, err := wayland.NewClient()
	if err != nil {
		return nil, err
	}

	result := &Client{Client: client}
	result.display.client = result
	result.display.id = result.NewObjectId()
	return result, nil
}

type Client struct {
	*wayland.Client
	display WlDisplay
}

func (client *Client) GetDisplay() WlDisplay {
	return client.display
}

type WlDisplay Object

func (object WlDisplay) Sync() WlCallback {
	callback := WlCallback(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, callback.id))

	return callback
}

func (object WlDisplay) GetRegistry() WlRegistry {
	registry := WlRegistry(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, registry.id))

	return registry
}

type WlRegistry Object

func (object WlRegistry) Bind(name uint32, iface string, version uint32) Object {
	id := Object{
		client: object.client,
		id: object.client.NewObjectId(),
	}

	object.client.Write(wayland.NewMessage(object.id, 0, name, iface, version, id.id))

	return id
}

type WlCallback Object

type WlCompositor Object

func (object WlCompositor) CreateSurface() WlSurface {
	id := WlSurface(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, id.id))

	return id
}

func (object WlCompositor) CreateRegion() WlRegion {
	id := WlRegion(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id))

	return id
}

type WlShmPool Object

func (object WlShmPool) CreateBuffer(offset int32, width int32, height int32, stride int32, format uint32) WlBuffer {
	id := WlBuffer(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, id.id, offset, width, height, stride, format))

	return id
}

func (object WlShmPool) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

func (object WlShmPool) Resize(size int32) {
	object.client.Write(wayland.NewMessage(object.id, 2, size))
}

type WlShm Object

func (object WlShm) CreatePool(fd int, size int32) WlShmPool {
	id := WlShmPool(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, id.id, size), fd)

	return id
}

func (object WlShm) Release() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type WlBuffer Object

func (object WlBuffer) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type WlDataOffer Object

func (object WlDataOffer) Accept(serial uint32, mimeType string) {
	object.client.Write(wayland.NewMessage(object.id, 0, serial, mimeType))
}

func (object WlDataOffer) Receive(mimeType string, fd int) {
	object.client.Write(wayland.NewMessage(object.id, 1, mimeType), fd)
}

func (object WlDataOffer) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

func (object WlDataOffer) Finish() {
	object.client.Write(wayland.NewMessage(object.id, 3))
}

func (object WlDataOffer) SetActions(dndActions uint32, preferredAction uint32) {
	object.client.Write(wayland.NewMessage(object.id, 4, dndActions, preferredAction))
}

type WlDataSource Object

func (object WlDataSource) Offer(mimeType string) {
	object.client.Write(wayland.NewMessage(object.id, 0, mimeType))
}

func (object WlDataSource) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

func (object WlDataSource) SetActions(dndActions uint32) {
	object.client.Write(wayland.NewMessage(object.id, 2, dndActions))
}

type WlDataDevice Object

func (object WlDataDevice) StartDrag(source WlDataSource, origin WlSurface, icon WlSurface, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, source.id, origin.id, icon.id, serial))
}

func (object WlDataDevice) SetSelection(source WlDataSource, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, source.id, serial))
}

func (object WlDataDevice) Release() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

type WlDataDeviceManager Object

func (object WlDataDeviceManager) CreateDataSource() WlDataSource {
	id := WlDataSource(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, id.id))

	return id
}

func (object WlDataDeviceManager) GetDataDevice(seat WlSeat) WlDataDevice {
	id := WlDataDevice(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, seat.id))

	return id
}

type WlShell Object

func (object WlShell) GetShellSurface(surface WlSurface) WlShellSurface {
	id := WlShellSurface(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, id.id, surface.id))

	return id
}

type WlShellSurface Object

func (object WlShellSurface) Pong(serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, serial))
}

func (object WlShellSurface) Move(seat WlSeat, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, seat.id, serial))
}

func (object WlShellSurface) Resize(seat WlSeat, serial uint32, edges uint32) {
	object.client.Write(wayland.NewMessage(object.id, 2, seat.id, serial, edges))
}

func (object WlShellSurface) SetToplevel() {
	object.client.Write(wayland.NewMessage(object.id, 3))
}

func (object WlShellSurface) SetTransient(parent WlSurface, x int32, y int32, flags uint32) {
	object.client.Write(wayland.NewMessage(object.id, 4, parent.id, x, y, flags))
}

func (object WlShellSurface) SetFullscreen(method uint32, framerate uint32, output WlOutput) {
	object.client.Write(wayland.NewMessage(object.id, 5, method, framerate, output.id))
}

func (object WlShellSurface) SetPopup(seat WlSeat, serial uint32, parent WlSurface, x int32, y int32, flags uint32) {
	object.client.Write(wayland.NewMessage(object.id, 6, seat.id, serial, parent.id, x, y, flags))
}

func (object WlShellSurface) SetMaximized(output WlOutput) {
	object.client.Write(wayland.NewMessage(object.id, 7, output.id))
}

func (object WlShellSurface) SetTitle(title string) {
	object.client.Write(wayland.NewMessage(object.id, 8, title))
}

func (object WlShellSurface) SetClass(class string) {
	object.client.Write(wayland.NewMessage(object.id, 9, class))
}

type WlSurface Object

func (object WlSurface) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WlSurface) Attach(buffer WlBuffer, x int32, y int32) {
	object.client.Write(wayland.NewMessage(object.id, 1, buffer.id, x, y))
}

func (object WlSurface) Damage(x int32, y int32, width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 2, x, y, width, height))
}

func (object WlSurface) Frame() WlCallback {
	callback := WlCallback(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 3, callback.id))

	return callback
}

func (object WlSurface) SetOpaqueRegion(region WlRegion) {
	object.client.Write(wayland.NewMessage(object.id, 4, region.id))
}

func (object WlSurface) SetInputRegion(region WlRegion) {
	object.client.Write(wayland.NewMessage(object.id, 5, region.id))
}

func (object WlSurface) Commit() {
	object.client.Write(wayland.NewMessage(object.id, 6))
}

func (object WlSurface) SetBufferTransform(transform int32) {
	object.client.Write(wayland.NewMessage(object.id, 7, transform))
}

func (object WlSurface) SetBufferScale(scale int32) {
	object.client.Write(wayland.NewMessage(object.id, 8, scale))
}

func (object WlSurface) DamageBuffer(x int32, y int32, width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 9, x, y, width, height))
}

func (object WlSurface) Offset(x int32, y int32) {
	object.client.Write(wayland.NewMessage(object.id, 10, x, y))
}

type WlSeat Object

func (object WlSeat) GetPointer() WlPointer {
	id := WlPointer(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, id.id))

	return id
}

func (object WlSeat) GetKeyboard() WlKeyboard {
	id := WlKeyboard(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id))

	return id
}

func (object WlSeat) GetTouch() WlTouch {
	id := WlTouch(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, id.id))

	return id
}

func (object WlSeat) Release() {
	object.client.Write(wayland.NewMessage(object.id, 3))
}

type WlPointer Object

func (object WlPointer) SetCursor(serial uint32, surface WlSurface, hotspotX int32, hotspotY int32) {
	object.client.Write(wayland.NewMessage(object.id, 0, serial, surface.id, hotspotX, hotspotY))
}

func (object WlPointer) Release() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type WlKeyboard Object

func (object WlKeyboard) Release() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type WlTouch Object

func (object WlTouch) Release() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type WlOutput Object

func (object WlOutput) Release() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type WlRegion Object

func (object WlRegion) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WlRegion) Add(x int32, y int32, width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 1, x, y, width, height))
}

func (object WlRegion) Subtract(x int32, y int32, width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 2, x, y, width, height))
}

type WlSubcompositor Object

func (object WlSubcompositor) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WlSubcompositor) GetSubsurface(surface WlSurface, parent WlSurface) WlSubsurface {
	id := WlSubsurface(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id, parent.id))

	return id
}

type WlSubsurface Object

func (object WlSubsurface) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WlSubsurface) SetPosition(x int32, y int32) {
	object.client.Write(wayland.NewMessage(object.id, 1, x, y))
}

func (object WlSubsurface) PlaceAbove(sibling WlSurface) {
	object.client.Write(wayland.NewMessage(object.id, 2, sibling.id))
}

func (object WlSubsurface) PlaceBelow(sibling WlSurface) {
	object.client.Write(wayland.NewMessage(object.id, 3, sibling.id))
}

func (object WlSubsurface) SetSync() {
	object.client.Write(wayland.NewMessage(object.id, 4))
}

func (object WlSubsurface) SetDesync() {
	object.client.Write(wayland.NewMessage(object.id, 5))
}

type WlFixes Object

func (object WlFixes) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WlFixes) DestroyRegistry(registry WlRegistry) {
	object.client.Write(wayland.NewMessage(object.id, 1, registry.id))
}

type ZwpLinuxDmabufV1 Object

func (object ZwpLinuxDmabufV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ZwpLinuxDmabufV1) CreateParams() ZwpLinuxBufferParamsV1 {
	paramsId := ZwpLinuxBufferParamsV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, paramsId.id))

	return paramsId
}

func (object ZwpLinuxDmabufV1) GetDefaultFeedback() ZwpLinuxDmabufFeedbackV1 {
	id := ZwpLinuxDmabufFeedbackV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, id.id))

	return id
}

func (object ZwpLinuxDmabufV1) GetSurfaceFeedback(surface WlSurface) ZwpLinuxDmabufFeedbackV1 {
	id := ZwpLinuxDmabufFeedbackV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 3, id.id, surface.id))

	return id
}

type ZwpLinuxBufferParamsV1 Object

func (object ZwpLinuxBufferParamsV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ZwpLinuxBufferParamsV1) Add(fd int, planeIdx uint32, offset uint32, stride uint32, modifierHi uint32, modifierLo uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, planeIdx, offset, stride, modifierHi, modifierLo), fd)
}

func (object ZwpLinuxBufferParamsV1) Create(width int32, height int32, format uint32, flags uint32) {
	object.client.Write(wayland.NewMessage(object.id, 2, width, height, format, flags))
}

func (object ZwpLinuxBufferParamsV1) CreateImmed(width int32, height int32, format uint32, flags uint32) WlBuffer {
	bufferId := WlBuffer(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 3, bufferId.id, width, height, format, flags))

	return bufferId
}

type ZwpLinuxDmabufFeedbackV1 Object

func (object ZwpLinuxDmabufFeedbackV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type WpPresentation Object

func (object WpPresentation) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpPresentation) Feedback(surface WlSurface) WpPresentationFeedback {
	callback := WpPresentationFeedback(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, surface.id, callback.id))

	return callback
}

type WpPresentationFeedback Object

type ZwpTabletManagerV2 Object

func (object ZwpTabletManagerV2) GetTabletSeat(seat WlSeat) ZwpTabletSeatV2 {
	tabletSeat := ZwpTabletSeatV2(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, tabletSeat.id, seat.id))

	return tabletSeat
}

func (object ZwpTabletManagerV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ZwpTabletSeatV2 Object

func (object ZwpTabletSeatV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type ZwpTabletToolV2 Object

func (object ZwpTabletToolV2) SetCursor(serial uint32, surface WlSurface, hotspotX int32, hotspotY int32) {
	object.client.Write(wayland.NewMessage(object.id, 0, serial, surface.id, hotspotX, hotspotY))
}

func (object ZwpTabletToolV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ZwpTabletV2 Object

func (object ZwpTabletV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type ZwpTabletPadRingV2 Object

func (object ZwpTabletPadRingV2) SetFeedback(description string, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, description, serial))
}

func (object ZwpTabletPadRingV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ZwpTabletPadStripV2 Object

func (object ZwpTabletPadStripV2) SetFeedback(description string, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, description, serial))
}

func (object ZwpTabletPadStripV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ZwpTabletPadGroupV2 Object

func (object ZwpTabletPadGroupV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type ZwpTabletPadV2 Object

func (object ZwpTabletPadV2) SetFeedback(button uint32, description string, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, button, description, serial))
}

func (object ZwpTabletPadV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ZwpTabletPadDialV2 Object

func (object ZwpTabletPadDialV2) SetFeedback(description string, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, description, serial))
}

func (object ZwpTabletPadDialV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type WpViewporter Object

func (object WpViewporter) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpViewporter) GetViewport(surface WlSurface) WpViewport {
	id := WpViewport(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type WpViewport Object

func (object WpViewport) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpViewport) SetSource(x wayland.Fixed, y wayland.Fixed, width wayland.Fixed, height wayland.Fixed) {
	object.client.Write(wayland.NewMessage(object.id, 1, x, y, width, height))
}

func (object WpViewport) SetDestination(width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 2, width, height))
}

type XdgWmBase Object

func (object XdgWmBase) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgWmBase) CreatePositioner() XdgPositioner {
	id := XdgPositioner(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id))

	return id
}

func (object XdgWmBase) GetXdgSurface(surface WlSurface) XdgSurface {
	id := XdgSurface(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, id.id, surface.id))

	return id
}

func (object XdgWmBase) Pong(serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 3, serial))
}

type XdgPositioner Object

func (object XdgPositioner) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgPositioner) SetSize(width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 1, width, height))
}

func (object XdgPositioner) SetAnchorRect(x int32, y int32, width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 2, x, y, width, height))
}

func (object XdgPositioner) SetAnchor(anchor uint32) {
	object.client.Write(wayland.NewMessage(object.id, 3, anchor))
}

func (object XdgPositioner) SetGravity(gravity uint32) {
	object.client.Write(wayland.NewMessage(object.id, 4, gravity))
}

func (object XdgPositioner) SetConstraintAdjustment(constraintAdjustment uint32) {
	object.client.Write(wayland.NewMessage(object.id, 5, constraintAdjustment))
}

func (object XdgPositioner) SetOffset(x int32, y int32) {
	object.client.Write(wayland.NewMessage(object.id, 6, x, y))
}

func (object XdgPositioner) SetReactive() {
	object.client.Write(wayland.NewMessage(object.id, 7))
}

func (object XdgPositioner) SetParentSize(parentWidth int32, parentHeight int32) {
	object.client.Write(wayland.NewMessage(object.id, 8, parentWidth, parentHeight))
}

func (object XdgPositioner) SetParentConfigure(serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 9, serial))
}

type XdgSurface Object

func (object XdgSurface) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgSurface) GetToplevel() XdgToplevel {
	id := XdgToplevel(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id))

	return id
}

func (object XdgSurface) GetPopup(parent XdgSurface, positioner XdgPositioner) XdgPopup {
	id := XdgPopup(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, id.id, parent.id, positioner.id))

	return id
}

func (object XdgSurface) SetWindowGeometry(x int32, y int32, width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 3, x, y, width, height))
}

func (object XdgSurface) AckConfigure(serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 4, serial))
}

type XdgToplevel Object

func (object XdgToplevel) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgToplevel) SetParent(parent XdgToplevel) {
	object.client.Write(wayland.NewMessage(object.id, 1, parent.id))
}

func (object XdgToplevel) SetTitle(title string) {
	object.client.Write(wayland.NewMessage(object.id, 2, title))
}

func (object XdgToplevel) SetAppId(appId string) {
	object.client.Write(wayland.NewMessage(object.id, 3, appId))
}

func (object XdgToplevel) ShowWindowMenu(seat WlSeat, serial uint32, x int32, y int32) {
	object.client.Write(wayland.NewMessage(object.id, 4, seat.id, serial, x, y))
}

func (object XdgToplevel) Move(seat WlSeat, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 5, seat.id, serial))
}

func (object XdgToplevel) Resize(seat WlSeat, serial uint32, edges uint32) {
	object.client.Write(wayland.NewMessage(object.id, 6, seat.id, serial, edges))
}

func (object XdgToplevel) SetMaxSize(width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 7, width, height))
}

func (object XdgToplevel) SetMinSize(width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 8, width, height))
}

func (object XdgToplevel) SetMaximized() {
	object.client.Write(wayland.NewMessage(object.id, 9))
}

func (object XdgToplevel) UnsetMaximized() {
	object.client.Write(wayland.NewMessage(object.id, 10))
}

func (object XdgToplevel) SetFullscreen(output WlOutput) {
	object.client.Write(wayland.NewMessage(object.id, 11, output.id))
}

func (object XdgToplevel) UnsetFullscreen() {
	object.client.Write(wayland.NewMessage(object.id, 12))
}

func (object XdgToplevel) SetMinimized() {
	object.client.Write(wayland.NewMessage(object.id, 13))
}

type XdgPopup Object

func (object XdgPopup) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgPopup) Grab(seat WlSeat, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, seat.id, serial))
}

func (object XdgPopup) Reposition(positioner XdgPositioner, token uint32) {
	object.client.Write(wayland.NewMessage(object.id, 2, positioner.id, token))
}

type WpAlphaModifierV1 Object

func (object WpAlphaModifierV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpAlphaModifierV1) GetSurface(surface WlSurface) WpAlphaModifierSurfaceV1 {
	id := WpAlphaModifierSurfaceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type WpAlphaModifierSurfaceV1 Object

func (object WpAlphaModifierSurfaceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpAlphaModifierSurfaceV1) SetMultiplier(factor uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, factor))
}

type WpColorManagerV1 Object

func (object WpColorManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpColorManagerV1) GetOutput(output WlOutput) WpColorManagementOutputV1 {
	id := WpColorManagementOutputV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, output.id))

	return id
}

func (object WpColorManagerV1) GetSurface(surface WlSurface) WpColorManagementSurfaceV1 {
	id := WpColorManagementSurfaceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, id.id, surface.id))

	return id
}

func (object WpColorManagerV1) GetSurfaceFeedback(surface WlSurface) WpColorManagementSurfaceFeedbackV1 {
	id := WpColorManagementSurfaceFeedbackV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 3, id.id, surface.id))

	return id
}

func (object WpColorManagerV1) CreateIccCreator() WpImageDescriptionCreatorIccV1 {
	obj := WpImageDescriptionCreatorIccV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 4, obj.id))

	return obj
}

func (object WpColorManagerV1) CreateParametricCreator() WpImageDescriptionCreatorParamsV1 {
	obj := WpImageDescriptionCreatorParamsV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 5, obj.id))

	return obj
}

func (object WpColorManagerV1) CreateWindowsScrgb() WpImageDescriptionV1 {
	imageDescription := WpImageDescriptionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 6, imageDescription.id))

	return imageDescription
}

type WpColorManagementOutputV1 Object

func (object WpColorManagementOutputV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpColorManagementOutputV1) GetImageDescription() WpImageDescriptionV1 {
	imageDescription := WpImageDescriptionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, imageDescription.id))

	return imageDescription
}

type WpColorManagementSurfaceV1 Object

func (object WpColorManagementSurfaceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpColorManagementSurfaceV1) SetImageDescription(imageDescription WpImageDescriptionV1, renderIntent uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, imageDescription.id, renderIntent))
}

func (object WpColorManagementSurfaceV1) UnsetImageDescription() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

type WpColorManagementSurfaceFeedbackV1 Object

func (object WpColorManagementSurfaceFeedbackV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpColorManagementSurfaceFeedbackV1) GetPreferred() WpImageDescriptionV1 {
	imageDescription := WpImageDescriptionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, imageDescription.id))

	return imageDescription
}

func (object WpColorManagementSurfaceFeedbackV1) GetPreferredParametric() WpImageDescriptionV1 {
	imageDescription := WpImageDescriptionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, imageDescription.id))

	return imageDescription
}

type WpImageDescriptionCreatorIccV1 Object

func (object WpImageDescriptionCreatorIccV1) Create() WpImageDescriptionV1 {
	imageDescription := WpImageDescriptionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, imageDescription.id))

	return imageDescription
}

func (object WpImageDescriptionCreatorIccV1) SetIccFile(iccProfile int, offset uint32, length uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, offset, length), iccProfile)
}

type WpImageDescriptionCreatorParamsV1 Object

func (object WpImageDescriptionCreatorParamsV1) Create() WpImageDescriptionV1 {
	imageDescription := WpImageDescriptionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, imageDescription.id))

	return imageDescription
}

func (object WpImageDescriptionCreatorParamsV1) SetTfNamed(tf uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, tf))
}

func (object WpImageDescriptionCreatorParamsV1) SetTfPower(eexp uint32) {
	object.client.Write(wayland.NewMessage(object.id, 2, eexp))
}

func (object WpImageDescriptionCreatorParamsV1) SetPrimariesNamed(primaries uint32) {
	object.client.Write(wayland.NewMessage(object.id, 3, primaries))
}

func (object WpImageDescriptionCreatorParamsV1) SetPrimaries(rX int32, rY int32, gX int32, gY int32, bX int32, bY int32, wX int32, wY int32) {
	object.client.Write(wayland.NewMessage(object.id, 4, rX, rY, gX, gY, bX, bY, wX, wY))
}

func (object WpImageDescriptionCreatorParamsV1) SetLuminances(minLum uint32, maxLum uint32, referenceLum uint32) {
	object.client.Write(wayland.NewMessage(object.id, 5, minLum, maxLum, referenceLum))
}

func (object WpImageDescriptionCreatorParamsV1) SetMasteringDisplayPrimaries(rX int32, rY int32, gX int32, gY int32, bX int32, bY int32, wX int32, wY int32) {
	object.client.Write(wayland.NewMessage(object.id, 6, rX, rY, gX, gY, bX, bY, wX, wY))
}

func (object WpImageDescriptionCreatorParamsV1) SetMasteringLuminance(minLum uint32, maxLum uint32) {
	object.client.Write(wayland.NewMessage(object.id, 7, minLum, maxLum))
}

func (object WpImageDescriptionCreatorParamsV1) SetMaxCll(maxCll uint32) {
	object.client.Write(wayland.NewMessage(object.id, 8, maxCll))
}

func (object WpImageDescriptionCreatorParamsV1) SetMaxFall(maxFall uint32) {
	object.client.Write(wayland.NewMessage(object.id, 9, maxFall))
}

type WpImageDescriptionV1 Object

func (object WpImageDescriptionV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpImageDescriptionV1) GetInformation() WpImageDescriptionInfoV1 {
	information := WpImageDescriptionInfoV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, information.id))

	return information
}

type WpImageDescriptionInfoV1 Object

type WpColorRepresentationManagerV1 Object

func (object WpColorRepresentationManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpColorRepresentationManagerV1) GetSurface(surface WlSurface) WpColorRepresentationSurfaceV1 {
	id := WpColorRepresentationSurfaceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type WpColorRepresentationSurfaceV1 Object

func (object WpColorRepresentationSurfaceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpColorRepresentationSurfaceV1) SetAlphaMode(alphaMode uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, alphaMode))
}

func (object WpColorRepresentationSurfaceV1) SetCoefficientsAndRange(coefficients uint32, rnge uint32) {
	object.client.Write(wayland.NewMessage(object.id, 2, coefficients, rnge))
}

func (object WpColorRepresentationSurfaceV1) SetChromaLocation(chromaLocation uint32) {
	object.client.Write(wayland.NewMessage(object.id, 3, chromaLocation))
}

type WpCommitTimingManagerV1 Object

func (object WpCommitTimingManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpCommitTimingManagerV1) GetTimer(surface WlSurface) WpCommitTimerV1 {
	id := WpCommitTimerV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type WpCommitTimerV1 Object

func (object WpCommitTimerV1) SetTimestamp(tvSecHi uint32, tvSecLo uint32, tvNsec uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, tvSecHi, tvSecLo, tvNsec))
}

func (object WpCommitTimerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type WpContentTypeManagerV1 Object

func (object WpContentTypeManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpContentTypeManagerV1) GetSurfaceContentType(surface WlSurface) WpContentTypeV1 {
	id := WpContentTypeV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type WpContentTypeV1 Object

func (object WpContentTypeV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpContentTypeV1) SetContentType(contentType uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, contentType))
}

type WpCursorShapeManagerV1 Object

func (object WpCursorShapeManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpCursorShapeManagerV1) GetPointer(pointer WlPointer) WpCursorShapeDeviceV1 {
	cursorShapeDevice := WpCursorShapeDeviceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, cursorShapeDevice.id, pointer.id))

	return cursorShapeDevice
}

func (object WpCursorShapeManagerV1) GetTabletToolV2(tabletTool ZwpTabletToolV2) WpCursorShapeDeviceV1 {
	cursorShapeDevice := WpCursorShapeDeviceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, cursorShapeDevice.id, tabletTool.id))

	return cursorShapeDevice
}

type WpCursorShapeDeviceV1 Object

func (object WpCursorShapeDeviceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpCursorShapeDeviceV1) SetShape(serial uint32, shape uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, serial, shape))
}

type WpDrmLeaseDeviceV1 Object

func (object WpDrmLeaseDeviceV1) CreateLeaseRequest() WpDrmLeaseRequestV1 {
	id := WpDrmLeaseRequestV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, id.id))

	return id
}

func (object WpDrmLeaseDeviceV1) Release() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type WpDrmLeaseConnectorV1 Object

func (object WpDrmLeaseConnectorV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type WpDrmLeaseRequestV1 Object

func (object WpDrmLeaseRequestV1) RequestConnector(connector WpDrmLeaseConnectorV1) {
	object.client.Write(wayland.NewMessage(object.id, 0, connector.id))
}

func (object WpDrmLeaseRequestV1) Submit() WpDrmLeaseV1 {
	id := WpDrmLeaseV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id))

	return id
}

type WpDrmLeaseV1 Object

func (object WpDrmLeaseV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type ExtBackgroundEffectManagerV1 Object

func (object ExtBackgroundEffectManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtBackgroundEffectManagerV1) GetBackgroundEffect(surface WlSurface) ExtBackgroundEffectSurfaceV1 {
	id := ExtBackgroundEffectSurfaceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type ExtBackgroundEffectSurfaceV1 Object

func (object ExtBackgroundEffectSurfaceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtBackgroundEffectSurfaceV1) SetBlurRegion(region WlRegion) {
	object.client.Write(wayland.NewMessage(object.id, 1, region.id))
}

type ExtDataControlManagerV1 Object

func (object ExtDataControlManagerV1) CreateDataSource() ExtDataControlSourceV1 {
	id := ExtDataControlSourceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, id.id))

	return id
}

func (object ExtDataControlManagerV1) GetDataDevice(seat WlSeat) ExtDataControlDeviceV1 {
	id := ExtDataControlDeviceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, seat.id))

	return id
}

func (object ExtDataControlManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

type ExtDataControlDeviceV1 Object

func (object ExtDataControlDeviceV1) SetSelection(source ExtDataControlSourceV1) {
	object.client.Write(wayland.NewMessage(object.id, 0, source.id))
}

func (object ExtDataControlDeviceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

func (object ExtDataControlDeviceV1) SetPrimarySelection(source ExtDataControlSourceV1) {
	object.client.Write(wayland.NewMessage(object.id, 2, source.id))
}

type ExtDataControlSourceV1 Object

func (object ExtDataControlSourceV1) Offer(mimeType string) {
	object.client.Write(wayland.NewMessage(object.id, 0, mimeType))
}

func (object ExtDataControlSourceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtDataControlOfferV1 Object

func (object ExtDataControlOfferV1) Receive(mimeType string, fd int) {
	object.client.Write(wayland.NewMessage(object.id, 0, mimeType), fd)
}

func (object ExtDataControlOfferV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtForeignToplevelListV1 Object

func (object ExtForeignToplevelListV1) Stop() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtForeignToplevelListV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtForeignToplevelHandleV1 Object

func (object ExtForeignToplevelHandleV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type ExtIdleNotifierV1 Object

func (object ExtIdleNotifierV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtIdleNotifierV1) GetIdleNotification(timeout uint32, seat WlSeat) ExtIdleNotificationV1 {
	id := ExtIdleNotificationV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, timeout, seat.id))

	return id
}

func (object ExtIdleNotifierV1) GetInputIdleNotification(timeout uint32, seat WlSeat) ExtIdleNotificationV1 {
	id := ExtIdleNotificationV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, id.id, timeout, seat.id))

	return id
}

type ExtIdleNotificationV1 Object

func (object ExtIdleNotificationV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type ExtImageCaptureSourceV1 Object

func (object ExtImageCaptureSourceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type ExtOutputImageCaptureSourceManagerV1 Object

func (object ExtOutputImageCaptureSourceManagerV1) CreateSource(output WlOutput) ExtImageCaptureSourceV1 {
	source := ExtImageCaptureSourceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, source.id, output.id))

	return source
}

func (object ExtOutputImageCaptureSourceManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtForeignToplevelImageCaptureSourceManagerV1 Object

func (object ExtForeignToplevelImageCaptureSourceManagerV1) CreateSource(toplevelHandle ExtForeignToplevelHandleV1) ExtImageCaptureSourceV1 {
	source := ExtImageCaptureSourceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, source.id, toplevelHandle.id))

	return source
}

func (object ExtForeignToplevelImageCaptureSourceManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtImageCopyCaptureManagerV1 Object

func (object ExtImageCopyCaptureManagerV1) CreateSession(source ExtImageCaptureSourceV1, options uint32) ExtImageCopyCaptureSessionV1 {
	session := ExtImageCopyCaptureSessionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, session.id, source.id, options))

	return session
}

func (object ExtImageCopyCaptureManagerV1) CreatePointerCursorSession(source ExtImageCaptureSourceV1, pointer WlPointer) ExtImageCopyCaptureCursorSessionV1 {
	session := ExtImageCopyCaptureCursorSessionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, session.id, source.id, pointer.id))

	return session
}

func (object ExtImageCopyCaptureManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

type ExtImageCopyCaptureSessionV1 Object

func (object ExtImageCopyCaptureSessionV1) CreateFrame() ExtImageCopyCaptureFrameV1 {
	frame := ExtImageCopyCaptureFrameV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, frame.id))

	return frame
}

func (object ExtImageCopyCaptureSessionV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtImageCopyCaptureFrameV1 Object

func (object ExtImageCopyCaptureFrameV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtImageCopyCaptureFrameV1) AttachBuffer(buffer WlBuffer) {
	object.client.Write(wayland.NewMessage(object.id, 1, buffer.id))
}

func (object ExtImageCopyCaptureFrameV1) DamageBuffer(x int32, y int32, width int32, height int32) {
	object.client.Write(wayland.NewMessage(object.id, 2, x, y, width, height))
}

func (object ExtImageCopyCaptureFrameV1) Capture() {
	object.client.Write(wayland.NewMessage(object.id, 3))
}

type ExtImageCopyCaptureCursorSessionV1 Object

func (object ExtImageCopyCaptureCursorSessionV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtImageCopyCaptureCursorSessionV1) GetCaptureSession() ExtImageCopyCaptureSessionV1 {
	session := ExtImageCopyCaptureSessionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, session.id))

	return session
}

type ExtSessionLockManagerV1 Object

func (object ExtSessionLockManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtSessionLockManagerV1) Lock() ExtSessionLockV1 {
	id := ExtSessionLockV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id))

	return id
}

type ExtSessionLockV1 Object

func (object ExtSessionLockV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtSessionLockV1) GetLockSurface(surface WlSurface, output WlOutput) ExtSessionLockSurfaceV1 {
	id := ExtSessionLockSurfaceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id, output.id))

	return id
}

func (object ExtSessionLockV1) UnlockAndDestroy() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

type ExtSessionLockSurfaceV1 Object

func (object ExtSessionLockSurfaceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtSessionLockSurfaceV1) AckConfigure(serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, serial))
}

type ExtTransientSeatManagerV1 Object

func (object ExtTransientSeatManagerV1) Create() ExtTransientSeatV1 {
	seat := ExtTransientSeatV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, seat.id))

	return seat
}

func (object ExtTransientSeatManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtTransientSeatV1 Object

func (object ExtTransientSeatV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type ExtWorkspaceManagerV1 Object

func (object ExtWorkspaceManagerV1) Commit() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtWorkspaceManagerV1) Stop() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtWorkspaceGroupHandleV1 Object

func (object ExtWorkspaceGroupHandleV1) CreateWorkspace(workspace string) {
	object.client.Write(wayland.NewMessage(object.id, 0, workspace))
}

func (object ExtWorkspaceGroupHandleV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type ExtWorkspaceHandleV1 Object

func (object ExtWorkspaceHandleV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object ExtWorkspaceHandleV1) Activate() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

func (object ExtWorkspaceHandleV1) Deactivate() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

func (object ExtWorkspaceHandleV1) Assign(workspaceGroup ExtWorkspaceGroupHandleV1) {
	object.client.Write(wayland.NewMessage(object.id, 3, workspaceGroup.id))
}

func (object ExtWorkspaceHandleV1) Remove() {
	object.client.Write(wayland.NewMessage(object.id, 4))
}

type WpFifoManagerV1 Object

func (object WpFifoManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpFifoManagerV1) GetFifo(surface WlSurface) WpFifoV1 {
	id := WpFifoV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type WpFifoV1 Object

func (object WpFifoV1) SetBarrier() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpFifoV1) WaitBarrier() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

func (object WpFifoV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

type WpFractionalScaleManagerV1 Object

func (object WpFractionalScaleManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpFractionalScaleManagerV1) GetFractionalScale(surface WlSurface) WpFractionalScaleV1 {
	id := WpFractionalScaleV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type WpFractionalScaleV1 Object

func (object WpFractionalScaleV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type WpLinuxDrmSyncobjManagerV1 Object

func (object WpLinuxDrmSyncobjManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpLinuxDrmSyncobjManagerV1) GetSurface(surface WlSurface) WpLinuxDrmSyncobjSurfaceV1 {
	id := WpLinuxDrmSyncobjSurfaceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

func (object WpLinuxDrmSyncobjManagerV1) ImportTimeline(fd int) WpLinuxDrmSyncobjTimelineV1 {
	id := WpLinuxDrmSyncobjTimelineV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, id.id), fd)

	return id
}

type WpLinuxDrmSyncobjTimelineV1 Object

func (object WpLinuxDrmSyncobjTimelineV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

type WpLinuxDrmSyncobjSurfaceV1 Object

func (object WpLinuxDrmSyncobjSurfaceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpLinuxDrmSyncobjSurfaceV1) SetAcquirePoint(timeline WpLinuxDrmSyncobjTimelineV1, pointHi uint32, pointLo uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, timeline.id, pointHi, pointLo))
}

func (object WpLinuxDrmSyncobjSurfaceV1) SetReleasePoint(timeline WpLinuxDrmSyncobjTimelineV1, pointHi uint32, pointLo uint32) {
	object.client.Write(wayland.NewMessage(object.id, 2, timeline.id, pointHi, pointLo))
}

type WpPointerWarpV1 Object

func (object WpPointerWarpV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpPointerWarpV1) WarpPointer(surface WlSurface, pointer WlPointer, x wayland.Fixed, y wayland.Fixed, serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 1, surface.id, pointer.id, x, y, serial))
}

type WpSecurityContextManagerV1 Object

func (object WpSecurityContextManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpSecurityContextManagerV1) CreateListener(listenFd int, closeFd int) WpSecurityContextV1 {
	id := WpSecurityContextV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id), listenFd, closeFd)

	return id
}

type WpSecurityContextV1 Object

func (object WpSecurityContextV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpSecurityContextV1) SetSandboxEngine(name string) {
	object.client.Write(wayland.NewMessage(object.id, 1, name))
}

func (object WpSecurityContextV1) SetAppId(appId string) {
	object.client.Write(wayland.NewMessage(object.id, 2, appId))
}

func (object WpSecurityContextV1) SetInstanceId(instanceId string) {
	object.client.Write(wayland.NewMessage(object.id, 3, instanceId))
}

func (object WpSecurityContextV1) Commit() {
	object.client.Write(wayland.NewMessage(object.id, 4))
}

type WpSinglePixelBufferManagerV1 Object

func (object WpSinglePixelBufferManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpSinglePixelBufferManagerV1) CreateU32RgbaBuffer(r uint32, g uint32, b uint32, a uint32) WlBuffer {
	id := WlBuffer(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, r, g, b, a))

	return id
}

type WpTearingControlManagerV1 Object

func (object WpTearingControlManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object WpTearingControlManagerV1) GetTearingControl(surface WlSurface) WpTearingControlV1 {
	id := WpTearingControlV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type WpTearingControlV1 Object

func (object WpTearingControlV1) SetPresentationHint(hint uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, hint))
}

func (object WpTearingControlV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type XdgActivationV1 Object

func (object XdgActivationV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgActivationV1) GetActivationToken() XdgActivationTokenV1 {
	id := XdgActivationTokenV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id))

	return id
}

func (object XdgActivationV1) Activate(token string, surface WlSurface) {
	object.client.Write(wayland.NewMessage(object.id, 2, token, surface.id))
}

type XdgActivationTokenV1 Object

func (object XdgActivationTokenV1) SetSerial(serial uint32, seat WlSeat) {
	object.client.Write(wayland.NewMessage(object.id, 0, serial, seat.id))
}

func (object XdgActivationTokenV1) SetAppId(appId string) {
	object.client.Write(wayland.NewMessage(object.id, 1, appId))
}

func (object XdgActivationTokenV1) SetSurface(surface WlSurface) {
	object.client.Write(wayland.NewMessage(object.id, 2, surface.id))
}

func (object XdgActivationTokenV1) Commit() {
	object.client.Write(wayland.NewMessage(object.id, 3))
}

func (object XdgActivationTokenV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 4))
}

type XdgWmDialogV1 Object

func (object XdgWmDialogV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgWmDialogV1) GetXdgDialog(toplevel XdgToplevel) XdgDialogV1 {
	id := XdgDialogV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, toplevel.id))

	return id
}

type XdgDialogV1 Object

func (object XdgDialogV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgDialogV1) SetModal() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

func (object XdgDialogV1) UnsetModal() {
	object.client.Write(wayland.NewMessage(object.id, 2))
}

type XdgSystemBellV1 Object

func (object XdgSystemBellV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgSystemBellV1) Ring(surface WlSurface) {
	object.client.Write(wayland.NewMessage(object.id, 1, surface.id))
}

type XdgToplevelDragManagerV1 Object

func (object XdgToplevelDragManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgToplevelDragManagerV1) GetXdgToplevelDrag(dataSource WlDataSource) XdgToplevelDragV1 {
	id := XdgToplevelDragV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, dataSource.id))

	return id
}

type XdgToplevelDragV1 Object

func (object XdgToplevelDragV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgToplevelDragV1) Attach(toplevel XdgToplevel, xOffset int32, yOffset int32) {
	object.client.Write(wayland.NewMessage(object.id, 1, toplevel.id, xOffset, yOffset))
}

type XdgToplevelIconManagerV1 Object

func (object XdgToplevelIconManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgToplevelIconManagerV1) CreateIcon() XdgToplevelIconV1 {
	id := XdgToplevelIconV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id))

	return id
}

func (object XdgToplevelIconManagerV1) SetIcon(toplevel XdgToplevel, icon XdgToplevelIconV1) {
	object.client.Write(wayland.NewMessage(object.id, 2, toplevel.id, icon.id))
}

type XdgToplevelIconV1 Object

func (object XdgToplevelIconV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgToplevelIconV1) SetName(iconName string) {
	object.client.Write(wayland.NewMessage(object.id, 1, iconName))
}

func (object XdgToplevelIconV1) AddBuffer(buffer WlBuffer, scale int32) {
	object.client.Write(wayland.NewMessage(object.id, 2, buffer.id, scale))
}

type XdgToplevelTagManagerV1 Object

func (object XdgToplevelTagManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XdgToplevelTagManagerV1) SetToplevelTag(toplevel XdgToplevel, tag string) {
	object.client.Write(wayland.NewMessage(object.id, 1, toplevel.id, tag))
}

func (object XdgToplevelTagManagerV1) SetToplevelDescription(toplevel XdgToplevel, description string) {
	object.client.Write(wayland.NewMessage(object.id, 2, toplevel.id, description))
}

type XwaylandShellV1 Object

func (object XwaylandShellV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XwaylandShellV1) GetXwaylandSurface(surface WlSurface) XwaylandSurfaceV1 {
	id := XwaylandSurfaceV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, surface.id))

	return id
}

type XwaylandSurfaceV1 Object

func (object XwaylandSurfaceV1) SetSerial(serialLo uint32, serialHi uint32) {
	object.client.Write(wayland.NewMessage(object.id, 0, serialLo, serialHi))
}

func (object XwaylandSurfaceV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type XxInputMethodV1 Object

func (object XxInputMethodV1) CommitString(text string) {
	object.client.Write(wayland.NewMessage(object.id, 0, text))
}

func (object XxInputMethodV1) SetPreeditString(text string, cursorBegin int32, cursorEnd int32) {
	object.client.Write(wayland.NewMessage(object.id, 1, text, cursorBegin, cursorEnd))
}

func (object XxInputMethodV1) DeleteSurroundingText(beforeLength uint32, afterLength uint32) {
	object.client.Write(wayland.NewMessage(object.id, 2, beforeLength, afterLength))
}

func (object XxInputMethodV1) Commit(serial uint32) {
	object.client.Write(wayland.NewMessage(object.id, 3, serial))
}

func (object XxInputMethodV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 4))
}

type XxInputMethodManagerV2 Object

func (object XxInputMethodManagerV2) GetInputMethod(seat WlSeat) XxInputMethodV1 {
	inputMethod := XxInputMethodV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 0, seat.id, inputMethod.id))

	return inputMethod
}

func (object XxInputMethodManagerV2) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

type XxSessionManagerV1 Object

func (object XxSessionManagerV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XxSessionManagerV1) GetSession(reason uint32, session string) XxSessionV1 {
	id := XxSessionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 1, id.id, reason, session))

	return id
}

type XxSessionV1 Object

func (object XxSessionV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XxSessionV1) Remove() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

func (object XxSessionV1) AddToplevel(toplevel XdgToplevel, name string) XxToplevelSessionV1 {
	id := XxToplevelSessionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 2, id.id, toplevel.id, name))

	return id
}

func (object XxSessionV1) RestoreToplevel(toplevel XdgToplevel, name string) XxToplevelSessionV1 {
	id := XxToplevelSessionV1(Object{
		client: object.client,
		id: object.client.NewObjectId(),
	})

	object.client.Write(wayland.NewMessage(object.id, 3, id.id, toplevel.id, name))

	return id
}

type XxToplevelSessionV1 Object

func (object XxToplevelSessionV1) Destroy() {
	object.client.Write(wayland.NewMessage(object.id, 0))
}

func (object XxToplevelSessionV1) Remove() {
	object.client.Write(wayland.NewMessage(object.id, 1))
}

