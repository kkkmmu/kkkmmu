1. bcmsdk可以被编译为user mode和kernel mode，并且根据最新的官方文档，以后将不会再支持kernel mode.
2. bcmsdk的makefile结构与web页面的布局有异曲同工之妙。
     1. 将编译子目录需要执行的指令放置到同一个文件中(Make.subdirs)，然后再需要的地方include该文件。
     2. 将编译lib的通用指令放置在同一文件中(Make.lib)，然后再需要编译为lib的Makefile中include该文件，并指定lib名称。
     3. 将全局的配置选项放置到同一个文件（Make.config）中，然后再需要的地方include该文件
   这样做的好处是避免了在各级目录Makefile中的重复指令, 使得构建过程易于维护。

3. User Mode:
	kernel_modules:
		1. $(MAKE) -C $(SDK)/systems/linux/kernel/modules kernel_version=$(kernel_version) OPT_CFLAGS="$(ADD_TO_CFLAGS)" subdirs="shared" override-target=linux-$(platform); (libkern)
		2. $(MAKE) -C $(SDK)/systems/linux/kernel/modules kernel_version=$(kernel_version) OPT_CFLAGS="$(ADD_TO_CFLAGS)" subdirs="uk-proxy" override-target=linux-$(platform); (linux-uk-proxy.ko)
		3. $(MAKE) -C $(SDK)/systems/bde/linux/kernel kernel_version=$(kernel_version) OPT_CFLAGS="$(ADD_TO_CFLAGS)";  (linux-kernel-bde.ko)
		4. $(MAKE) -C $(SDK)/systems/bde/linux/user/kernel kernel_version=$(kernel_version) OPT_CFLAGS="$(ADD_TO_CFLAGS)"  (linux-user-bde.ko)
	user_libs:
		5. $(MAKE) -C $(SDK)/systems/bde/linux/user CFLAGS="$(CFLAGS)" (liblubde)
		6. $(MAKE) -C $(SDK)/src CFLAGS="$(ADD_TO_CFLAGS)"  (libSUBDIR)
		7. $(MAKE) -C $(SDK)/systems/drv CFLAGS="$(CFLAGS)"  (libdrivers)
	bcm.user: (User mode bcmshell)
		8. $(OBJCOPY) --strip-debug bcm.user bcm.user.dbg

4. Kernel Mode: 
	kernel_modules:
		1. $(MAKE) -C $(SDK)/src  (libSUBDIR)
		2. $(MAKE) -C $(SDK)/systems/drv (libdrivers)
		3. $(MAKE) -C $(SDK)/systems/bde/linux/kernel  (linux-kernel-bde.ko)
		4. $(MAKE) -C $(SDK)/systems/linux/kernel/modules kernel_version=$(kernel_version)      (For kernel modules: linux-uk-proxy.ko bcm-diag-full.ko, bcm-diag.ko, bcm-core.ko, bcm-net.ko)
	user_apps: 
		5. $(MAKE) -C $(SDK)/src/sal/appl (libsal_appl)
		6. $(MAKE) -C $(SDK)/systems/linux/kernel/user/bcm-diag-proxy BCM_PROXY=$(BCM_PROXY) override-target=$(user-override) bldroot_suffix=/$(override-target)     （FOR bcm.user.proxy）

5. 整个代码的编译入口:
	user mode: 	system/linux/user/PLATFORM/
	kernel mode: 	system/linux/kernel/PLATFORM/
   在该目录先执行make以后会进入Make.linux，在Make.linux中会根据LINUX_MAKE_USER是否设置连选择是进入systems/linux/user/common还是systems/linux/kernel/common进行下一步操作。
   在common目录中就会指定该模式（user/kernel）所需要的目标及编译规则。

6. 至于每个目录要编译什么，要编译成什么，则由该目录的Makefile以及Make.config, Make.targets控制。

7.如何根据需求对SDK进行裁剪:
Customers can build SDK with the necessary XGS chip support, PHY chip support, and other features. This can reduce the bcm.user image size if there is any concern for storage size. 

Makefiles
$SDK/make/Make.local and $SDK/make/make.config will determine which chip/PHY/features to build into bcm.user.

Modify the following files to get a smaller bcm.user image.

Steps
1. Copy $SDK/make/Make.local.template to $SDK/make/Make.local 

2. Modify Make.local  

FEATURE_LIST: FEATURE_LIST is the list of features to include. If you do not define FEATURE_LIST in your Make.local, Make.config will have a default value for different chip settings. Define FEATURE_LIST and remove unused features to get a smaller bcm.user image.

Use BCM_PTL_SPT to select the Strata switch chip needed. By default, the driver supports all Strata switch and fabric chips. Uncomment the line for BCM_PTL_SPT (partial support) and chips to support. For example, if you just want your SDK to support HR2, let BCM_PTL_SPT = 1 and BCM_56150_A0 = 1 in Make.local.

BCM_PHY_LIST is a list of PHYs to include. The default is to include all of them. Select the necessary PHY and put them in BCM_PHY_LIST.

BCMX APIs have not been enhanced or supported for newer devices since SDK-5.10.2. You should set INCLUDE_BCMX=0 to remove BCMX from SDK. 

Set DISPATCH_LIST as ESW for enterprise platforms only. Or set it as ROBO for robo only. 

3. Re-compile SDK to get a smaller bcm.user image for partial chip/PHY/feature support. 
