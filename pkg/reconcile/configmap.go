package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO 注释作用域是包吗，还是文件里的代码，还是说整个controller
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// ConfigMaps reconciles the configmap(s) required for the instance in the current context.
func ConfigMaps(ctx context.Context, instance client.Object, params models.Params) error {

	var obj ApolloObject

	// TODO switch 修改一下
	if instance.GetObjectKind().GroupVersionKind().Kind == "ApolloPortal" {
		obj = ApolloPortal()
	}

	desired := obj.DesiredConfigMaps(ctx, instance, params)

	// 可以把资源全抽象为接口，直接放reconcile下，reconcile下文件中策略模式，选择不同的apollo对象，文件夹下的文件可以按照 构建对象、创建更新、删除 的模式

	// TODO 可以优化为先获取create、upodate、delete列表，然后再统一apply

	// first, handle the create/update parts
	if err := obj.ExpectedConfigMaps(ctx, instance, params, desired, true); err != nil {
		return fmt.Errorf("failed to reconcile the expected configmaps: %w", err)
	}

	// then, delete the extra objects
	if err := obj.DeleteConfigMaps(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the configmaps to be deleted: %w", err)
	}

	return nil
}
